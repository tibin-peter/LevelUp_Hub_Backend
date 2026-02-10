package booking

import "time"

type CreateBookingRequest struct {
	SlotID uint `json:"slot_id"`
}

type BookingResponseDTO struct {
	BookingID   uint      `json:"booking_id"`
	StudentName string    `json:"student_name"`
	MentorName  string    `json:"mentor_name"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
}
