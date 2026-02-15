package payment

type BookingPort interface {
    GetBookingByID(id uint) (BookingDTO, error)
    MarkBookingPaid(id uint) error
}

type BookingDTO struct {
    ID        uint
    StudentID uint
		MentorID  uint
    Price     float64
    Status    string
}