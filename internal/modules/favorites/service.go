package favorites

import "errors"

type Service interface {
	AddFavorite(studentID, mentorProfileID uint) error 
	RemoveFavorite(studentID, mentorProfileID uint) error
	ListFavorites(studentID uint) ([]FavoriteResponse, error) 
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) AddFavorite(studentID, mentorProfileID uint) error {

	if mentorProfileID == 0 {
		return errors.New("invalid mentor profile")
	}

	exists, err := s.repo.Exists(studentID, mentorProfileID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already favorited")
	}

	fav := &Favorite{
		StudentID:       studentID,
		MentorProfileID: mentorProfileID,
	}

	return s.repo.Create(fav)
}

func (s *service) RemoveFavorite(studentID, mentorProfileID uint) error {
	return s.repo.Delete(studentID, mentorProfileID)
}

func (s *service) ListFavorites(studentID uint) ([]FavoriteResponse, error) {

	return s.repo.GetDetailedByStudent(studentID)
}