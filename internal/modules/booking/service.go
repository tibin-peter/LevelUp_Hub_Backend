package booking

import (
	"LevelUp_Hub_Backend/internal/modules/profile"
	"LevelUp_Hub_Backend/internal/modules/slot"
	"errors"
	"log"
)

type Service interface {
	CreateBooking(studentID, slotID uint) error
	GetStudentUpcoming(studentID uint) ([]BookingResponseDTO, error)
	GetStudentHistory(studentID uint) ([]BookingResponseDTO, error)
	CancelBooking(bookingID, studentID uint) error

	ApproveBooking(bookingID, mentorUserID uint) error
	RejectBooking(bookingID, mentorUserID uint) error
	GetMentorUpcoming(userID uint) ([]BookingResponseDTO, error)
	GetMentorHistory(userID uint) ([]BookingResponseDTO, error)

	

	GetStudentBookings(studentID uint) ([]Booking, error)
	GetMentorBookings(mentorUserID uint) ([]Booking, error)
}

type service struct {
	repo       Repository
	slotrepo   slot.Repository
	mentorrepo profile.MentorRepository
}

func NewService(repo Repository, s slot.Repository, m profile.MentorRepository) Service {
	return &service{
		repo:       repo,
		slotrepo:   s,
		mentorrepo: m,
	}
}

//create booking by student
func (s *service) CreateBooking(studentID, slotID uint) error {
	slotData, err := s.slotrepo.GetByID(slotID)
	if err != nil {
		return err
	}
	if slotData.IsBooked {
		return errors.New("slot already booked")
	}
	mentor, err := s.mentorrepo.FindMentorByID(slotData.MentorProfileID)
	if err != nil {
		return err
	}
	booking := Booking{
		StudentID:       studentID,
		MentorProfileID: slotData.MentorProfileID,
		SlotID:          slotID,
		Status:          "pending",
		Price:           mentor.HourlyPrice,
	}
	return s.repo.Create(&booking)
}

//approve booking by mentor
func (s *service) ApproveBooking(userID, mentorUserID uint) error {
	booking, err := s.repo.GetByID(userID)
	log.Println("service:",booking)
	if err != nil {
		return err
	}
	mentorProfileID, err := s.slotrepo.GetProfileIDByUserID(mentorUserID)
	if err != nil {
		return err
	}
	if booking.MentorProfileID != mentorProfileID {
		return errors.New("can't approve this booking")
	}
	if booking.Status != "pending" {
		return errors.New("invalid status")
	}
	//update status
	if err := s.repo.UpdateStatus(userID, "confirmed"); err != nil {
		return err
	}
	//lock the slot
	return s.slotrepo.MarkBooked(booking.SlotID, true)
}

//reject booking by mentor
func (s *service) RejectBooking(bookingID, mentorUserID uint) error {
	booking, err := s.repo.GetByID(bookingID)
	if err != nil {
		return err
	}

	mentorProfileID, _ := s.slotrepo.GetProfileIDByUserID(mentorUserID)
	if booking.MentorProfileID != mentorProfileID {
		return errors.New("can't reject this booking")
	}
	return s.repo.UpdateStatus(bookingID, "rejected")
}

//cancell booking by student
func (s *service) CancelBooking(bookingID, studentID uint) error {

	booking, err := s.repo.GetByID(bookingID)
	if err != nil {
		return err
	}

	if booking.StudentID != studentID {
		return errors.New("not your booking")
	}

	if booking.Status != "pending" && booking.Status != "confirmed" {
		return errors.New("cannot cancel")
	}

	// free slot if confirmed
	if booking.Status == "confirmed" {
		s.slotrepo.MarkBooked(booking.SlotID, false)
	}

	return s.repo.UpdateStatus(bookingID, "cancelled")
}

//get all booking for student
func (s *service) GetStudentBookings(studentID uint) ([]Booking, error) {
	return s.repo.GetStudentBookings(studentID)
}

//get all for mentor
func (s *service) GetMentorBookings(mentorUserID uint) ([]Booking, error) {

	mentorProfileID, err := s.slotrepo.GetProfileIDByUserID(mentorUserID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetMentorBookings(mentorProfileID)
}

func (s *service) GetStudentUpcoming(studentID uint) ([]BookingResponseDTO, error) {
	return s.repo.GetUpcomingByStudent(studentID)
}

func (s *service) GetStudentHistory(studentID uint) ([]BookingResponseDTO, error) {
	return s.repo.GetHistoryByStudent(studentID)
}

func (s *service) GetMentorUpcoming(userID uint) ([]BookingResponseDTO, error) {

	mpID, err := s.slotrepo.GetProfileIDByUserID(userID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetUpcomingByMentor(mpID)
}

func (s *service) GetMentorHistory(userID uint) ([]BookingResponseDTO, error) {

	mpID, err := s.slotrepo.GetProfileIDByUserID(userID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetHistoryByMentor(mpID)
}
