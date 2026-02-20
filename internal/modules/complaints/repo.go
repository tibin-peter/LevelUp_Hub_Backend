package complaints

import (
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(c *Complaint) error
	GetByUser(userID uint) ([]Complaint, error)
	GetAll(search string) ([]Complaint, error)

	UpdateReply(id uint, reply, status string) error
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

// create
func (r *repo) Create(c *Complaint) error {
	return r.db.Create(c).Error
}

// get own complaints
func (r *repo) GetByUser(userID uint) ([]Complaint, error) {
	var list []Complaint

	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&list).Error

	return list, err
}

// get all
func (r *repo) GetAll(search string) ([]Complaint, error) {

	var list []Complaint

	query := r.db.
		Order("created_at DESC")

	if search != "" {
		query = query.Where(
			"category ILIKE ? OR subject ILIKE ? OR description ILIKE ? OR status ILIKE ?",
			"%"+search+"%",
			"%"+search+"%",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	err := query.Find(&list).Error

	return list, err
}

// update and reply
func (r *repo) UpdateReply(id uint, reply string, status string) error {

	return r.db.
		Model(&Complaint{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"admin_reply": reply,
			"status":      status,
			"updated_at":  time.Now(),
		}).Error
}
