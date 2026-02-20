package usermanagement

import (
	"LevelUp_Hub_Backend/internal/modules/profile"

	"gorm.io/gorm"
)

type Repository interface {
	ListUsers(role, search string) ([]profile.User, error)
	BlockUser(id uint) error
	UnBlockUser(id uint) error
	PendingMentors() ([]profile.MentorProfile,error)
	ApproveMentor(id uint) error
	RejectMentor(id uint) error
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

///////////////// User /////////////////
func (r *repo) ListUsers(role, search string) ([]profile.User, error) {

	var users []profile.User

	query := r.db.Where("role = ?", role)

	if search != "" {
		query = query.Where(
			"name ILIKE ? OR email ILIKE ?",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	err := query.Find(&users).Error

	return users, err
}

func (r *repo) BlockUser(id uint) error {
	return r.db.
		Model(&profile.User{}).
		Where("id=? AND role != ?", id, "admin").
		Update("is_blocked", true).Error
}

func (r *repo) UnBlockUser(id uint) error {
	return r.db.
		Model(&profile.User{}).
		Where("id=?", id).
		Update("is_blocked", false).Error
}

///////////////// Mentor /////////////////
func (r *repo) PendingMentors() ([]profile.MentorProfile,error) {

	var list []profile.MentorProfile

	err := r.db.
		Preload("User").
		Where("status = ?", "pending").
		Find(&list).Error

	return list, err
}

func (r *repo) ApproveMentor(id uint) error {
	return r.db.
		Model(&profile.MentorProfile{}).
		Where("id=? AND status=?", id,"pending").
		Update("status", "approved").Error
}

func (r *repo) RejectMentor(id uint) error {
	return r.db.
		Model(&profile.MentorProfile{}).
		Where("id=?", id).
		Update("status", "rejected").Error
}
