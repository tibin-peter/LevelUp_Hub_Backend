package complaints

import "errors"

type Service interface {
	CreateComplaint(userID uint, role string, req CreateComplaintRequest) error
	MyComplaints(userID uint) ([]Complaint, error)
	AllComplaints() ([]Complaint, error)
	AdminReply(id uint, reply string, status string) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

// create complaints
func (s *service) CreateComplaint(userID uint, role string, req CreateComplaintRequest) error {

	c := &Complaint{
		UserID:      userID,
		Role:        role,
		Category:    req.Category,
		Subject:     req.Subject,
		Description: req.Description,
		Status:      "open",
	}

	return s.repo.Create(c)
}

// get complains
func (s *service) MyComplaints(userID uint) ([]Complaint, error) {
	return s.repo.GetByUser(userID)
}

// get all
func (s *service) AllComplaints() ([]Complaint, error) {
	return s.repo.GetAll()
}

// admin reply
func (s *service) AdminReply(id uint, reply string, status string) error {
	if reply == "" {
		return errors.New("reply required")
	}
	return s.repo.UpdateReply(id, reply, status)
}