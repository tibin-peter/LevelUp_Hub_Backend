package auth

import (
	"LevelUp_Hub_Backend/internal/modules/profile"
	"LevelUp_Hub_Backend/internal/utils"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service interface {
	SendOTP(email string) (error)
	Register(name, email, password, role,inputOTP string) (string,string,*profile.User,error)
	Login(email, password string) (string, string,*profile.User, error)
	Refresh(refresh string) (string, string, error)
	Logout(userID uint)error
}

type service struct {
	repo      profile.Repository
	mentorrepo  profile.MentorRepository
	redis     *redis.Client
	JWTSecret string
}

func NewService(repo profile.Repository,mentorrepo profile.MentorRepository, rdb *redis.Client, jwtSecret string) *service {
	return &service{
		repo:      repo,
		mentorrepo: mentorrepo,
		redis:     rdb,
		JWTSecret:jwtSecret,
	}
}

//func for send otp
func (s *service) SendOTP(email string) (error) {
	if email == "" {
		return fmt.Errorf("invalid email")
	}
	//generate otp
	otp, err := utils.GenerateOTP()
	if err != nil {
		return err
	}
	ctx := context.Background()
	//store otp
	err1 := s.redis.Set(
		ctx,
		"otp:"+email,
		otp,
		5*time.Minute,
	).Err()
	if err1 != nil {
		return err1
	}
	//send email
	if err := utils.SendOTPEmail(email, otp); err != nil {
		return err
	}
	log.Println(otp)

	return nil
}

//fucn for Register
func (s *service) Register(name, email, password, role,inputOTP string) (string,string,*profile.User,error) {
	ctx:=context.Background()

	//get otp
	storedotp, err := s.redis.Get(ctx, "otp:"+email).Result()
	if err != nil {
		return "","",nil,fmt.Errorf("otp expired")
	}
	//check
	if storedotp != inputOTP {
		return "","",nil,fmt.Errorf("invalid otp")
	}

	//check user exist
	existing, err := s.repo.FindUserByEmail(email)
	if err != nil {
    return "", "",nil,err // real DB error
}
	if existing.ID != 0{
		return "","",nil,fmt.Errorf("email already registed")
	}
	//validate role
	if role != "student" && role != "mentor" {
		return "","",nil,fmt.Errorf("invalid role")
	}
	//hasing password
	hash, err := utils.HashPassword(password)
	if err != nil {
		return "","",nil,err
	}
	//create user
	newUser := profile.User{
		Name:       name,
		Email:      email,
		Password:   hash,
		Role:       role,
		IsVerified: true,
	}
	//save user
	if err := s.repo.CreateUser(&newUser); err != nil {
		return "","",nil,err
	}

	// autoe create mentorprofile with userid
	if role == "mentor" {
		mp := profile.MentorProfile{
			UserID: newUser.ID,
		}
		if err := s.mentorrepo.CreateMentor(&mp); err != nil {
			return "", "", nil, err
		}
	}
	//delete existing otp in redis
	s.redis.Del(ctx,"otp:"+email)
	//generate access token
	access, _, err := utils.GenerateAccessToken(newUser.ID,newUser.Email,newUser.Role,s.JWTSecret)
	if err != nil {
		return "", "",nil, err
	}
	//generate refresh
	refresh, exp, err := utils.GenerateRefreshToken(newUser.ID,newUser.Email,newUser.Role,s.JWTSecret)
	if err != nil {
		return "", "",nil, err
	}

  //store refresh to redis
	s.redis.Set(
		ctx,
		fmt.Sprintf("refresh:%d", newUser.ID),
		refresh,
		time.Until(exp),
	)

	return access, refresh,&newUser, nil
}

//fucn login
func (s *service) Login(email, password string) (string, string, *profile.User, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return "", "",nil, err
	}
	if user.ID == 0 {
		return "", "",nil, fmt.Errorf("invalid email")
	}
	if !user.IsVerified {
		return "", "",nil, fmt.Errorf("user not varified")
	}
	if err := utils.CheckPassword(user.Password, password); err != nil {
		return "", "",nil, fmt.Errorf("wrong password")
	}
	access, _, err := utils.GenerateAccessToken(user.ID, user.Email, user.Role, s.JWTSecret)
	if err != nil {
		return "", "",nil, err
	}
	//generate refresh
	refresh, exp, err := utils.GenerateRefreshToken(user.ID, user.Email, user.Role, s.JWTSecret)
	if err != nil {
		return "", "",nil, err
	}
	//store in redis
	ctx := context.Background()
	s.redis.Set(
		ctx,
		fmt.Sprintf("refresh:%d", user.ID),
		refresh,
		time.Until(exp),
	)
	return access, refresh,user, nil

}

//func for refresh and token rotation
func (s *service) Refresh(refresh string) (string, string, error) {
	//valid jwt
	claims, err := utils.ValidateToken(refresh, s.JWTSecret)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}
	userID := claims.UserID
	ctx := context.Background()

	//get stored token from redis
	stored, err := s.redis.Get(
		ctx,
		fmt.Sprintf("refresh:%d", userID),
	).Result()
	if err != nil {
		return "", "", fmt.Errorf("refresh expired")
	}

	//compare
	if stored != refresh {
		return "", "", fmt.Errorf("refresh not matching")
	}

	//delete old for rotation
	s.redis.Del(ctx, fmt.Sprintf("refresh:%d", userID))

	//enerate new access
	newAccess, _, err := utils.GenerateAccessToken(
		userID,
		claims.Email,
		claims.Role,
		s.JWTSecret,
	)
	if err != nil {
		return "", "", err
	}

	//generate new refresh
	newRefresh, exp, err := utils.GenerateRefreshToken(
		userID,
		claims.Email,
		claims.Role,
		s.JWTSecret,
	)
	if err != nil {
		return "", "", err
	}

	// store new refresh
	s.redis.Set(
		ctx,
		fmt.Sprintf("refresh:%d", userID),
		newRefresh,
		time.Until(exp),
	)

	return newAccess, newRefresh, nil
}

//func for logout
func(s *service)Logout(userID uint)error{
	ctx:=context.Background()
	err:=s.redis.Del(
		ctx,
		fmt.Sprintf("refresh:%d",userID),
	).Err()
	if err!=nil{
		return err
	}
	return nil
}