package adminprofile

import (
	"LevelUp_Hub_Backend/internal/modules/profile"
	"LevelUp_Hub_Backend/internal/utils"
	"errors"
)

type Service interface {
	GetAdminProfile(id uint)(*profile.User,error)
	UpdateProfile(id uint,updated UpdateProfile)(*profile.User,error)
	UpdateProfilePicture(id uint,updated UpdateProfilePicture) (*profile.User,error)
	ChangePassword(id uint,input ChangePassword)error
}

type service struct {
	userRepo profile.Repository
}

func NewService(u profile.Repository) Service {
	return &service{
		userRepo: u,
	}
}

func(s *service)GetAdminProfile(id uint)(*profile.User,error){
	profile,err:=s.userRepo.FindUserById(id)
	return profile,err
}

func (s *service)UpdateProfile(id uint,updated UpdateProfile)(*profile.User,error){
	profile,err:=s.userRepo.FindUserById(id)
	if err!=nil{
		return nil,err
	}

	profile.Name=updated.Name
	profile.Email=updated.Email

	err1:= s.userRepo.UpdateUser(profile)
	if err1!=nil{
		return nil,err1
	}
	return profile,nil
}
func (s *service) UpdateProfilePicture(id uint,updated UpdateProfilePicture) (*profile.User,error) {
	profile,err:=s.userRepo.FindUserById(id)
	if err!=nil{
		return nil,err
	}
	profile.ProfilePicURL=updated.ProfilePicURL
	if err:=s.userRepo.UpdateUser(profile);err!=nil{
		return nil,err
	}
	return profile,nil
}

func (s *service)ChangePassword(id uint,input ChangePassword)error{
	user,err:=s.userRepo.FindUserById(id)
	if err!=nil{
		return err
	}
	if err:=utils.CheckPassword(user.Password,input.OldPassword);err!=nil{
		return errors.New("old password is not correct")
	}
	if input.NewPassword!=input.ConformPassword{
		return errors.New("newpassword and conoform password is not matching")
	}
	hash,err:=utils.HashPassword(input.NewPassword)
	if err!=nil{
		return err
	}
	user.Password=hash
	return s.userRepo.UpdateUser(user)
}