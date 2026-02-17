package booking

import (
	"LevelUp_Hub_Backend/internal/repository/generic"
	"fmt"
	// "time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(booking *Booking) error
	GetBySlotID(slotID uint) (*Booking, error)

	GetByID(id uint) (*Booking, error)
	UpdateStatus(id uint, status string) error

	MarkBookingPaid(id uint) error

	GetStudentBookings(studentID uint) ([]Booking, error)
	GetMentorBookings(mentorProfileID uint) ([]Booking, error)

	GetUpcomingByStudent(studentID uint) ([]BookingResponseDTO, error)
	GetHistoryByStudent(studentID uint) ([]BookingResponseDTO, error)

	GetUpcomingByMentor(mentorProfileID uint) ([]BookingResponseDTO, error)
	GetHistoryByMentor(mentorProfileID uint) ([]BookingResponseDTO, error)

	// GetStalePendingBookings(expiry time.Time) ([]Booking, error)
	// GetUnapprovedPastBookings(now time.Time) ([]Booking, error)
}

type repo struct {
	base *generic.Repository[Booking]
	db   *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{
		base: generic.NewRepository[Booking](db),
		db:   db,
	}
}

// create
func (r *repo) Create(b *Booking) error {
	return r.db.Create(b).Error
}

// get by slot id
func (r *repo) GetBySlotID(slotID uint) (*Booking, error) {
	var b Booking
	err := r.db.Where("slot_id = ?", slotID).First(&b).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// get by id
func (r *repo) GetByID(id uint) (*Booking, error) {
	var b Booking
	err := r.db.First(&b, id).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// update
func (r *repo) UpdateStatus(id uint, status string) error {
	err := r.db.Model(&Booking{}).Where("id = ?", id).Update("status", status).Error
	return err
}

func (r *repo) MarkBookingPaid(id uint) error {
	return r.db.Model(&Booking{}).
		Where("id = ?", id).
		Update("status", "paid").Error
}

// get bookings for student
func (r *repo) GetStudentBookings(studentID uint) ([]Booking, error) {
	var list []Booking
	err := r.db.Where("student_id = ?", studentID).Find(&list).Error
	return list, err
}

// get bookings for mentor
func (r *repo) GetMentorBookings(mentorID uint) ([]Booking, error) {
	var list []Booking
	err := r.db.Where("mentor_profile_id = ?", mentorID).Find(&list).Error
	return list, err
}

// upcoming
func (r *repo) GetUpcomingByStudent(studentID uint) ([]BookingResponseDTO, error) {

	var list []BookingResponseDTO

	err := r.db.Table("bookings").
		Select(`
			bookings.id as booking_id,
			s.name as student_name,
			mu.name as mentor_name,
			ms.start_time,
			ms.end_time,
			bookings.price,
			bookings.status
		`).
		Joins("JOIN users s ON s.id = bookings.student_id").
		Joins("JOIN mentor_profiles mp ON mp.id = bookings.mentor_profile_id").
		Joins("JOIN users mu ON mu.id = mp.user_id").
		Joins("JOIN mentor_slots ms ON ms.id = bookings.slot_id").
		Where(`
			bookings.student_id = ?
			AND bookings.status IN ('confirmed', 'paid')
			AND ms.start_time > NOW()
		`, studentID).
		Order("ms.start_time ASC").
		Scan(&list).Error
	fmt.Println("repo data:", list)

	return list, err
}

// for history
func (r *repo) GetHistoryByStudent(studentID uint) ([]BookingResponseDTO, error) {

	var list []BookingResponseDTO

	err := r.db.Table("bookings").
		Select(`
			bookings.id as booking_id,
			s.name as student_name,
			mu.name as mentor_name,
			ms.start_time,
			ms.end_time,
			bookings.price,
			bookings.status
		`).
		Joins("JOIN users s ON s.id = bookings.student_id").
		Joins("JOIN mentor_profiles mp ON mp.id = bookings.mentor_profile_id").
		Joins("JOIN users mu ON mu.id = mp.user_id").
		Joins("JOIN mentor_slots ms ON ms.id = bookings.slot_id").
		Where(`
			bookings.student_id = ?
			AND (
				ms.start_time < NOW()
				OR bookings.status IN ('completed','cancelled','rejected')
			)
		`, studentID).
		Order("ms.start_time DESC").
		Scan(&list).Error

	return list, err
}

func (r *repo) GetUpcomingByMentor(mid uint) ([]BookingResponseDTO, error) {

	var list []BookingResponseDTO

	err := r.db.Table("bookings").
		Select(`
			bookings.id as booking_id,
			s.name as student_name,
			mu.name as mentor_name,
			ms.start_time,
			ms.end_time,
			bookings.price,
			bookings.status
		`).
		Joins("JOIN users s ON s.id = bookings.student_id").
		Joins("JOIN mentor_profiles mp ON mp.id = bookings.mentor_profile_id").
		Joins("JOIN users mu ON mu.id = mp.user_id").
		Joins("JOIN mentor_slots ms ON ms.id = bookings.slot_id").
		Where(`
			bookings.mentor_profile_id = ?
			AND bookings.status IN ('confirmed', 'paid')
			AND ms.start_time > NOW()
		`, mid).
		Order("ms.start_time ASC").
		Scan(&list).Error

	fmt.Printf("[DEBUG] repo.GetUpcomingByMentor for mid %d returning %d items: %+v\n", mid, len(list), list)

	return list, err
}

func (r *repo) GetHistoryByMentor(mid uint) ([]BookingResponseDTO, error) {

	var list []BookingResponseDTO

	err := r.db.Table("bookings").
		Select(`
			bookings.id as booking_id,
			s.name as student_name,
			mu.name as mentor_name,
			ms.start_time,
			ms.end_time,
			bookings.price,
			bookings.status
		`).
		Joins("JOIN users s ON s.id = bookings.student_id").
		Joins("JOIN mentor_profiles mp ON mp.id = bookings.mentor_profile_id").
		Joins("JOIN users mu ON mu.id = mp.user_id").
		Joins("JOIN mentor_slots ms ON ms.id = bookings.slot_id").
		Where(`
			bookings.mentor_profile_id = ?
			AND (
				ms.start_time < NOW()
				OR bookings.status IN ('completed','cancelled','rejected')
			)
		`, mid).
		Order("ms.start_time DESC").
		Scan(&list).Error

	return list, err
}
