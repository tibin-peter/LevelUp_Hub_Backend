package profile

import "fmt"

type Service interface {
	GetUserById(id uint) (*User, error)
	FindUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id uint) error

	//mentor related
	CrateMentorProfile(userID uint,input MentorProfileInput)(*MentorProfile, error)
	GetMentorProfileByUserID(userID uint) (*MentorProfile, error)
	UpdateMentorProfile(userID uint, input MentorProfileInput) (*MentorProfile, error)
	GetPublicMentorProfile(mentorID uint) (*MentorProfile, error)
}

type service struct {
	repo Repository
	mentorrepo MentorRepository
}

func NewService(repo Repository,mentorrepo MentorRepository) Service {
	return &service{
		repo: repo,
		mentorrepo:mentorrepo ,
	}
}

// func for findbyid
func (s *service) GetUserById(id uint) (*User, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid user id")
	}
	user,err:=s.repo.FindUserById(id)
	if err!=nil{
		return nil,fmt.Errorf("user not found")
	}
	return user,nil
}

// func for findbyemail
func (s *service) FindUserByEmail(email string) (*User, error) {
	if email==""{
		return nil,fmt.Errorf("email is required")
	}
		user,err:=s.repo.FindUserByEmail(email)
		if err!=nil{
			return nil,fmt.Errorf("user not found")
		}
	return user,nil
}

// for update
func (s *service) UpdateUser(user *User) error {
	if user.ID==0{
		return fmt.Errorf("invalid user id")
	}
	return s.repo.UpdateUser(user)
}

// for deleteuser
func (s *service) DeleteUser(id uint) error {
	if id==0{
		return fmt.Errorf("invalid user id")
	}
	if err:=s.repo.DeleteUser(id);err!=nil{
		return fmt.Errorf("user not found")
	}
	return nil
}
//create mentor profile
func(s *service)CrateMentorProfile(userID uint,input MentorProfileInput)(*MentorProfile,error){
	// check user
	user, err := s.repo.FindUserById(userID)
	if err != nil {
		return nil,err
	}

	// ensure role is mentor
	if user.Role != "mentor" {
		return nil,fmt.Errorf("user is not a mentor")
	}

	// check existing profile
	existing, _ := s.mentorrepo.FindMentorByID(userID)
	if existing != nil {
		return nil,fmt.Errorf("mentor profile already exists")
	}

	profile := MentorProfile{
		UserID:          userID,
		Bio:             input.Bio,
		Skills:          input.Skills,
		Languages:       input.Languages,
		HourlyPrice:     input.HourlyPrice,
		ExperienceYears: input.ExperienceYears,
	}
	err1:=s.mentorrepo.CreateMentor(&profile)

	return &profile,err1
}

//get mentor profile by id
func (s *service) GetMentorProfileByUserID(userID uint) (*MentorProfile, error) {
	profile, err := s.mentorrepo.FindMentorByID(userID)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

//update mentor profile
func (s *service) UpdateMentorProfile(userID uint, input MentorProfileInput) (*MentorProfile, error) {
	profile, err := s.mentorrepo.FindMentorByID(userID)
	if err != nil {
		return nil,err
	}

	profile.Bio = input.Bio
	profile.Skills = input.Skills
	profile.Languages = input.Languages
	profile.HourlyPrice = input.HourlyPrice
	profile.ExperienceYears = input.ExperienceYears

	err1:=s.mentorrepo.UpdateMentor(profile)

	return profile,err1
}

//getpublic mentor profile
func (s *service) GetPublicMentorProfile(mentorID uint) (*MentorProfile, error) {
	profile, err := s.mentorrepo.FindMentorByID(mentorID)
	if err != nil {
		return nil, err
	}

	return profile, nil
}