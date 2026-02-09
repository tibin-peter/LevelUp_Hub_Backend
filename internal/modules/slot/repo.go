package slot

import (
	"LevelUp_Hub_Backend/internal/repository/generic"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	CreateSlot(slot *MentorSlot) error
	GetSlotsByMentor(mentorID uint) ([]MentorSlot, error)
	GetAvailableSlotsByMentor(mentoID uint) ([]MentorSlot, error)
	DeleteSlot(slotID uint, mentorID uint) error
	HasOverlap(mentorID uint, start time.Time, end time.Time) (bool, error)
	GetSlotsByDate(mentorID uint,date time.Time,) ([]MentorSlot, error)
	GetProfileIDByUserID(userID uint) (uint, error)
}

type repo struct {
	base *generic.Repository[MentorSlot]
	db   *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{
		base: generic.NewRepository[MentorSlot](db),
		db:   db,
	}
}

//create slot
func (r *repo) CreateSlot(slot *MentorSlot) error {
	return r.base.Create(slot)
}

//get slots by mentor
func (r *repo) GetSlotsByMentor(mentorID uint) ([]MentorSlot, error) {
	var slots []MentorSlot
	err := r.db.Where("mentor_profile_id = ?", mentorID).
		Order("start_time ASC").Find(&slots).Error

	return slots, err
}

//get availabel slots by mentor
func (r *repo) GetAvailableSlotsByMentor(mentoID uint) ([]MentorSlot, error) {
	var slots []MentorSlot

	err := r.db.
		Where(`
	mentor_profile_id = ?
	AND is_booked = false
	AND start_time > NOW()
	`, mentoID).
		Order("start_time ASC").
		Find(&slots).Error

	return slots, err
}

//for delete slot
func (r *repo) DeleteSlot(slotID uint, mentorProfileID uint) error {
	return r.db.Where("id = ? AND mentor_profile_id = ?", slotID, mentorProfileID).
		Delete(&MentorSlot{}).Error
}

//for prevent overlap
func (r *repo) HasOverlap(mentorID uint, start time.Time, end time.Time) (bool, error) {
	var count int64

	err := r.db.Model(&MentorSlot{}).
		Where(`
	mentor_profile_id = ?
	AND start_time < ?
	AND end_time > ?
	`, mentorID, end, start).Count(&count).Error

	return count > 0, err
}

//get by date
func (r *repo) GetSlotsByDate(mentorID uint,date time.Time,) ([]MentorSlot, error) {

	var slots []MentorSlot

	err := r.db.
		Where(`
   mentor_profile_id = ?
   AND DATE(start_time) = DATE(?)
   AND is_booked = false
  `, mentorID, date).
		Find(&slots).Error

	return slots, err
}

func (r *repo) GetProfileIDByUserID(userID uint) (uint, error) {
	var profile struct {
		ID uint
	}

	err := r.db.
		Table("mentor_profiles").
		Select("id").
		Where("user_id = ?", userID).
		First(&profile).Error

	return profile.ID, err
}