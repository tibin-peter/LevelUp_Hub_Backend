package connections

import (
	"LevelUp_Hub_Backend/internal/modules/profile"
	"errors"
)

type Service interface {
	CreateConnection(studentID, mentorProfileID uint) error
	GetMyMentors(studentID uint) ([]ConnectedMentor, error)
	GetStudentCount(mentorProfileID uint) (int64, error)
	GetMyStudents(mentorProfileID uint) ([]ConnectedStudent, error)
}

type service struct {
	repo       Repository
	mentorRepo profile.MentorRepository
}

func NewService(r Repository,mrepo profile.MentorRepository) Service {
	return &service{
		repo: r,
		mentorRepo: mrepo,
		}
}

func (s *service) CreateConnection(studentID, mentorProfileID uint) error {

	exists, err := s.repo.Exists(studentID, mentorProfileID)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	conn := &Connection{
		StudentID:       studentID,
		MentorProfileID: mentorProfileID,
	}

	return s.repo.Create(conn)
}

func (s *service) GetMyMentors(studentID uint) ([]ConnectedMentor, error) {
	return s.repo.GetConnectedMentors(studentID)
}

func (s *service) GetStudentCount(userID uint) (int64, error) {
	mentorProfile, err := s.mentorRepo.FindMentorByUserID(userID)
	if err != nil {
		return 0, errors.New("mentor profile not found")
	}
	return s.repo.CountStudents(mentorProfile.ID)
}

func (s *service) GetMyStudents(userID uint) ([]ConnectedStudent, error) {
  mentorProfile, err := s.mentorRepo.FindMentorByUserID(userID)
	if err != nil {
		return nil, errors.New("mentor profile not found")
	}
	return s.repo.GetConnectedStudents(mentorProfile.ID)
}