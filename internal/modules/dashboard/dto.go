package dashboard

import "time"

type StudentDashboard struct {
	ActiveBookings    int64 `json:"active_bookings"`
	CompletedSessions int64 `json:"completed_sessions"`
	FavoriteMentors   int64 `json:"favorite_mentors"`

	UpcomingSessions []UpcomingSession `json:"upcoming_sessions"`
	ConnectedMentors []ConnectedMentor `json:"connected_mentors"`
}

type MentorDashboard struct {
	TotalEarnings  float64 `json:"total_earnings"`
	TotalStudents  int64   `json:"total_students"`
	AvgRating      float64 `json:"avg_rating"`
	BookingRequests int64  `json:"booking_requests"`
}

type UpcomingSession struct {
	BookingID   uint      `json:"booking_id"`
	MentorName  string    `json:"mentor_name"`
	StartTime time.Time `json:"start_time"`
	EndTime time.Time `json:"end_time"`
}

type ConnectedMentor struct {
	MentorProfileID uint    `json:"mentor_profile_id"`
	Name            string  `json:"name"`
	ProfilePicURL   string  `json:"profile_pic_url"`
	Category        string  `json:"category"`
	HourlyPrice     float64 `json:"hourly_price"`
	RatingAvg       float64 `json:"rating_avg"`
}