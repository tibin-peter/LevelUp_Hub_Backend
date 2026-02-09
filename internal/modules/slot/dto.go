package slot

import "time"

type CreateSlotRequest struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type SlotResponse struct {
	ID        uint      `json:"id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	IsBooked  bool      `json:"is_booked"`
}


