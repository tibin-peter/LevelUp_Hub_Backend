package complaints

type Service interface {
	CreateComplaint(userID uint, role string, req CreateComplaintRequest) error
	MyComplaints(userID uint) ([]Complaint, error)
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
