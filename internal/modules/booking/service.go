package booking

import (
	"LevelUp_Hub_Backend/internal/modules/payment"
	"LevelUp_Hub_Backend/internal/modules/profile"
	"LevelUp_Hub_Backend/internal/modules/slot"
	"errors"
	"log"
	"time"
)

type Service interface {
	CreateBooking(studentID, slotID uint) (*Booking, error)
	GetStudentUpcoming(studentID uint) ([]BookingResponseDTO, error)
	GetStudentHistory(studentID uint) ([]BookingResponseDTO, error)
	CancelBooking(bookingID, studentID uint) error

	MarkBookingPaid(bookingID uint) error
	GetBookingByID(id uint) (payment.BookingDTO, error)
	UpdateStatus(id uint, status string) error
	FreeSlot(slotID uint) error
	SetPaymentPort(p PaymentPort)

	ApproveBooking(bookingID, mentorUserID uint) error
	RejectBooking(bookingID, mentorUserID uint) error
	GetMentorUpcoming(userID uint) ([]BookingResponseDTO, error)
	GetMentorHistory(userID uint) ([]BookingResponseDTO, error)
	CheckAndCompletePastSessions(mentorProfileID uint) error

	GetStudentBookings(studentID uint) ([]Booking, error)
	GetMentorBookings(mentorUserID uint) ([]Booking, error)
}

type service struct {
	repo       Repository
	slotrepo   slot.Repository
	mentorrepo profile.MentorRepository
	paymentSvc PaymentPort
}

func NewService(repo Repository, s slot.Repository, m profile.MentorRepository, pay PaymentPort) Service {
	return &service{
		repo:       repo,
		slotrepo:   s,
		mentorrepo: m,
		paymentSvc: pay,
	}
}

func (s *service) SetPaymentPort(p PaymentPort) {
	s.paymentSvc = p
}

// create booking by student
func (s *service) CreateBooking(studentID, slotID uint) (*Booking, error) { // Return *Booking
	// 1. Check if slot is already fully booked/finalized
	slotData, err := s.slotrepo.GetByID(slotID)
	if err != nil {
		return nil, err
	}
	if slotData.IsBooked {
		return nil, errors.New("this slot is already finalized and booked")
	}

	// Industry Standard 1: Minimum 1-hour slot check (if not already enforced in slot creation)
	duration := slotData.EndTime.Sub(slotData.StartTime)
	if duration < time.Hour {
		return nil, errors.New("mentorship sessions must be at least 1 hour")
	}

	// Industry Standard 2: Prevent booking slots that start too soon
	if time.Until(slotData.StartTime) < time.Hour {
		return nil, errors.New("cannot book a session starting in less than 1 hour")
	}

	// 2. Check if this student (or anyone) already has a pending booking for this slot
	existing, _ := s.repo.GetBySlotID(slotID)
	if existing != nil {
		if existing.Status == "pending_payment" && existing.StudentID == studentID {
			// Return the existing booking so the frontend can proceed to payment
			return existing, nil
		}
		return nil, errors.New("this slot is currently held for payment by another user")
	}

	// 3. If no booking exists, create a new one
	mentor, err := s.mentorrepo.FindMentorByID(slotData.MentorProfileID)
	if err != nil {
		return nil, err
	}

	booking := Booking{
		StudentID:       studentID,
		MentorProfileID: slotData.MentorProfileID,
		SlotID:          slotID,
		Status:          "pending_payment",
		Price:           mentor.HourlyPrice,
	}

	err = s.repo.Create(&booking)
	return &booking, err
}

// approve booking by mentor
func (s *service) ApproveBooking(bookingID, mentorUserID uint) error {
	booking, err := s.repo.GetByID(bookingID)
	log.Println("service:", booking)
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
	if booking.Status != "paid" {
		return errors.New("booking not paid yet")
	}
	//update status
	if err := s.repo.UpdateStatus(bookingID, "confirmed"); err != nil {
		return err
	}

	return s.paymentSvc.ReleaseEscrow(bookingID)
}

// reject booking by mentor
func (s *service) RejectBooking(bookingID, mentorUserID uint) error {

	booking, err := s.repo.GetByID(bookingID)
	if err != nil {
		return err
	}

	mentorProfileID, _ := s.slotrepo.GetProfileIDByUserID(mentorUserID)

	if booking.MentorProfileID != mentorProfileID {
		return errors.New("not your booking")
	}

	if err := s.repo.UpdateStatus(bookingID, "rejected"); err != nil {
		return err
	}

	s.slotrepo.MarkBooked(booking.SlotID, false)

	return s.paymentSvc.RefundEscrow(bookingID)
}

// cancell booking by student
func (s *service) CancelBooking(bookingID, studentID uint) error {

	booking, err := s.repo.GetByID(bookingID)
	if err != nil {
		return err
	}

	if booking.StudentID != studentID {
		return errors.New("not your booking")
	}

	if booking.Status == "confirmed" {
		slot, _ := s.slotrepo.GetByID(booking.SlotID)
		// Rule 7: Only allow cancellation 1 hour before start
		if time.Until(slot.StartTime) < time.Hour {
			return errors.New("cancellations only allowed 1 hour before session start")
		}

		// Trigger refund via Payment Module
		if err := s.paymentSvc.RefundEscrow(bookingID); err != nil {
			return err
		}
	}

	// free slot if confirmed
	s.slotrepo.MarkBooked(booking.SlotID, false)

	return s.repo.UpdateStatus(bookingID, "cancelled")
}

// get all booking for student
func (s *service) GetStudentBookings(studentID uint) ([]Booking, error) {
	return s.repo.GetStudentBookings(studentID)
}

// get all for mentor
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

	s.CheckAndCompletePastSessions(mpID)

	return s.repo.GetUpcomingByMentor(mpID)
}

func (s *service) GetMentorHistory(userID uint) ([]BookingResponseDTO, error) {

	mpID, err := s.slotrepo.GetProfileIDByUserID(userID)
	if err != nil {
		return nil, err
	}

	s.CheckAndCompletePastSessions(mpID)

	return s.repo.GetHistoryByMentor(mpID)
}

func (s *service) CheckAndCompletePastSessions(mentorProfileID uint) error {
	// Find confirmed bookings in the past
	bookings, err := s.repo.GetMentorBookings(mentorProfileID)
	if err != nil {
		return err
	}

	for _, b := range bookings {
		if b.Status == "confirmed" {
			slot, err := s.slotrepo.GetByID(b.SlotID)
			if err != nil {
				continue
			}
			// If session end time is in the past
			if slot.EndTime.Before(time.Now()) {
				log.Printf("[DEBUG] Completing past session: Booking %d", b.ID)
				s.repo.UpdateStatus(b.ID, "completed")
				s.paymentSvc.ReleaseEscrow(b.ID)
			}
		}
	}
	return nil
}

func (s *service) UpdateStatus(id uint, status string) error {
	return s.repo.UpdateStatus(id, status)
}

func (s *service) FreeSlot(slotID uint) error {
	return s.slotrepo.MarkBooked(slotID, false)
}

func (s *service) MarkBookingPaid(bookingID uint) error {

	booking, err := s.repo.GetByID(bookingID)
	if err != nil || booking == nil {
		return errors.New("booking not found")
	}

	if booking.Status != "pending_payment" {
		return errors.New("invalid booking status")
	}

	// mark paid
	log.Printf("[DEBUG] Marking booking %d as paid", bookingID)
	if err := s.repo.UpdateStatus(bookingID, "paid"); err != nil {
		return err
	}

	log.Printf("[DEBUG] Locking slot %d", booking.SlotID)
	// lock slot
	return s.slotrepo.MarkBooked(booking.SlotID, true)
}

func (s *service) GetBookingByID(id uint) (payment.BookingDTO, error) {

	b, err := s.repo.GetByID(id)
	if err != nil {
		return payment.BookingDTO{}, err
	}

	mentor, err := s.mentorrepo.FindMentorByID(b.MentorProfileID)
	if err != nil {
		return payment.BookingDTO{}, err
	}

	return payment.BookingDTO{
		ID:           b.ID,
		StudentID:    b.StudentID,
		MentorUserID: mentor.UserID,
		Price:        b.Price,
		Status:       b.Status,
	}, nil
}
