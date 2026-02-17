package ratings

import "time"

type CreateRatingRequest struct {
	BookingID uint   `json:"booking_id"`
	Rating    int    `json:"rating"`
	Review    string `json:"review"`
}

type CreateRatingResponse struct {
	AvgRating    float64 `json:"avg_rating"`
	TotalReviews int64   `json:"total_reviews"`
}

type MentorRatingSummary struct {
	AvgRating    float64 `json:"avg_rating"`
	TotalReviews int64   `json:"total_reviews"`
}

type RatingResponse struct {
	StudentID uint      `json:"student_id"`
	Rating    int       `json:"rating"`
	Review    string    `json:"review"`
	CreatedAt time.Time `json:"created_at"`
}

type TopMentor struct {
	MentorID  uint    `json:"mentor_id"`
	AvgRating float64 `json:"avg_rating"`
	Reviews   int64   `json:"reviews"`
}