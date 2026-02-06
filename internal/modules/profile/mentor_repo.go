package profile

import (
	"LevelUp_Hub_Backend/internal/repository/generic"
	"errors"

	"gorm.io/gorm"
)

type MentorRepository interface {
	CreateMentor(profile *MentorProfile) error
	FindMentorByUserID(userID uint) (*MentorProfile, error)
	FindMentorByID(id uint) (*MentorProfile, error)
	UpdateMentor(profile *MentorProfile) error
}

type mentorRepo struct {
	base *generic.Repository[MentorProfile]
	db   *gorm.DB
}

func NewMentorRepository(db *gorm.DB)MentorRepository {
	return &mentorRepo{
		base: generic.NewRepository[MentorProfile](db),
		db:   db,
	}
}

func (r *mentorRepo) CreateMentor(profile *MentorProfile) error {
	return r.base.Create(profile)
}

func (r *mentorRepo) FindMentorByUserID(userID uint) (*MentorProfile, error) {
	var profile MentorProfile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &profile, err
}

func (r *mentorRepo) FindMentorByID(id uint) (*MentorProfile, error) {
	return r.base.FindById(id)
}

func (r *mentorRepo) UpdateMentor(profile *MentorProfile) error {
	return r.base.Update(profile)
}