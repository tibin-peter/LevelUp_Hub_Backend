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
	FindMentorByIDPublic(userID uint) (*MentorProfileResponse, error)
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

func (r *mentorRepo) FindMentorByIDPublic(mentorID uint) (*MentorProfileResponse, error) {

	var result MentorProfileResponse

	err := r.db.
		Table("mentor_profiles").
		Select(`
			mentor_profiles.id,
			mentor_profiles.user_id,
			users.name,
			mentor_profiles.bio,
			mentor_profiles.skills,
			mentor_profiles.languages,
			mentor_profiles.hourly_price,
			mentor_profiles.experience_years,
			mentor_profiles.rating_avg,
			mentor_profiles.total_reviews,
			mentor_profiles.created_at,
			mentor_profiles.updated_at
		`).
		Joins("JOIN users ON users.id = mentor_profiles.user_id").
		Where("mentor_profiles.id = ?", mentorID).
		First(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}