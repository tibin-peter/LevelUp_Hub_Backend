package usermanagement

import (
	"LevelUp_Hub_Backend/internal/modules/complaints"
	"LevelUp_Hub_Backend/internal/modules/profile"
	"LevelUp_Hub_Backend/internal/utils"
)

type Service interface {
	ListUsers(role, search string) ([]profile.User, error)
	BlockUser(id uint) error
	UnblockUser(id uint) error

	PendingMentors() ([]profile.MentorProfile, error)
	ApproveMentor(id uint) error
	RejectMentor(id uint) error

	ListComplaints(search string) ([]complaints.Complaint, error)
	ReplyComplaint(id uint, reply string, status string) error
	ListAllUsers(search string) ([]profile.User, error)
	CreateUser(name, email, password, role string) error
}

type service struct {
	repo          Repository
	complaintRepo complaints.Repository
}

func NewService(r Repository, c complaints.Repository) Service {
	return &service{
		repo:          r,
		complaintRepo: c,
	}
}

///////////////// UserRelated ///////////////
func (s *service) ListUsers(role, search string) ([]profile.User, error) {

	return s.repo.ListUsers(role, search)
}

func (s *service) BlockUser(id uint) error {
	return s.repo.BlockUser(id)
}

func (s *service) UnblockUser(id uint) error {
	return s.repo.UnBlockUser(id)
}

///////////////// MentorRelated ///////////////
func (s *service) PendingMentors() ([]profile.MentorProfile, error) {

	return s.repo.PendingMentors()
}

func (s *service) ApproveMentor(id uint) error {
	return s.repo.ApproveMentor(id)
}

func (s *service) RejectMentor(id uint) error {
	return s.repo.RejectMentor(id)
}

///////////////// ComplaintRelated ///////////////
func (s *service) ListComplaints(search string) ([]complaints.Complaint, error) {

	return s.complaintRepo.GetAll(search)
}

func (s *service) ReplyComplaint(id uint, reply string, status string) error {

	return s.complaintRepo.UpdateReply(id, reply, status)
}

func (s *service) ListAllUsers(search string) ([]profile.User, error) {
	return s.repo.ListAllUsers(search)
}

func (s *service) CreateUser(name, email, password, role string) error {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	user := profile.User{
		Name:       name,
		Email:      email,
		Password:   hash,
		Role:       role,
		IsVerified: true,
	}
	if err := s.repo.CreateUser(&user); err != nil {
		return err
	}
	if role == "mentor" {
		mp := profile.MentorProfile{
			UserID: user.ID,
			Status: "approved",
		}
		if err := s.repo.CreateMentorProfile(&mp); err != nil {
			return err
		}
	}
	return nil
}