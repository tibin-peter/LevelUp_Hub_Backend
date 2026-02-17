package payment

type BookingPort interface {
	GetBookingByID(id uint) (BookingDTO, error)
	MarkBookingPaid(id uint) error
	UpdateStatus(id uint, status string) error
	FreeSlot(slotID uint) error
}

type BookingDTO struct {
	ID           uint
	StudentID    uint
	MentorID     uint
	MentorUserID uint
	Price        float64
	Status       string
}
