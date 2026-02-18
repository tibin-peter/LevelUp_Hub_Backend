package connections

type ConnectedMentor struct {
	MentorProfileID uint    `json:"mentor_profile_id"`
	Name            string  `json:"name"`
	ProfilePicURL   string  `json:"profile_pic_url"`
	Category        string  `json:"category"`
	HourlyPrice     float64 `json:"hourly_price"`
	RatingAvg       float64 `json:"rating_avg"`
}

type CountResponse struct {
	Total int64 `json:"total"`
}

type ConnectedStudent struct {
	StudentID     uint   `json:"student_id"`
	Name          string `json:"name"`
	ProfilePicURL string `json:"profile_pic_url"`
}