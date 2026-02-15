package booking

import (
	"LevelUp_Hub_Backend/internal/modules/payment"
	"LevelUp_Hub_Backend/internal/modules/profile"
	"LevelUp_Hub_Backend/internal/modules/slot"
	"errors"
	"log"
)

type Service interface {
	CreateBooking(studentID, slotID uint) (*Booking, error)
	GetStudentUpcoming(studentID uint) ([]BookingResponseDTO, error)
	GetStudentHistory(studentID uint) ([]BookingResponseDTO, error)
	CancelBooking(bookingID, studentID uint) error
	
	MarkBookingPaid(bookingID uint) error
	GetBookingByID(id uint) (payment.BookingDTO, error)
	SetPaymentPort(p PaymentPort)

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
	paymentSvc PaymentPort
}

func NewService(repo Repository, s slot.Repository, m profile.MentorRepository,pay PaymentPort) Service {
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

//create booking by student
func (s *service) CreateBooking(studentID, slotID uint) (*Booking, error) { // Return *Booking
    // 1. Check if slot is already fully booked/finalized
    slotData, err := s.slotrepo.GetByID(slotID)
    if err != nil { return nil, err }
    if slotData.IsBooked {
        return nil, errors.New("this slot is already finalized and booked")
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
    if err != nil { return nil, err }

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

//approve booking by mentor
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

	return s.paymentSvc.RefundEscrow(bookingID)
}

//reject booking by mentor
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

//cancell booking by student
func (s *service) CancelBooking(bookingID, studentID uint) error {

	booking, err := s.repo.GetByID(bookingID)
	if err != nil {
		return err
	}

	if booking.StudentID != studentID {
		return errors.New("not your booking")
	}

	if booking.Status != "pending_payment" && booking.Status != "paid" {
		return errors.New("cannot cancel")
	}

	// free slot if confirmed
	if booking.Status == "paid" || booking.Status == "confirmed" {
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

func (s *service) MarkBookingPaid(bookingID uint) error {

	booking, err := s.repo.GetByID(bookingID)
	if err != nil || booking == nil {
		return errors.New("booking not found")
	}

	if booking.Status != "pending_payment" {
		return errors.New("invalid booking status")
	}

	// mark paid
	if err := s.repo.UpdateStatus(bookingID, "paid"); err != nil {
		return err
	}

	// lock slot
	return s.slotrepo.MarkBooked(booking.SlotID, true)
}


func (s *service) GetBookingByID(id uint) (payment.BookingDTO, error) {

	b, err := s.repo.GetByID(id)
	if err != nil {
		return payment.BookingDTO{}, err
	}

	return payment.BookingDTO{
		ID:        b.ID,
		StudentID: b.StudentID,
		MentorID:  b.MentorProfileID,
		Price:     b.Price,
		Status:    b.Status,
	}, nil
}