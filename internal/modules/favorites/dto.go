package favorites

type AddFavoriteRequest struct {
	MentorProfileID uint `json:"mentor_profile_id"`
}

type FavoriteResponse struct {
	MentorProfileID uint    `json:"mentor_profile_id"`
	Name            string  `json:"name"`
	ProfilePicURL   string  `json:"profile_pic_url"`
	Category        string  `json:"category"`
	HourlyPrice     float64 `json:"hourly_price"`
	RatingAvg       float64 `json:"rating_avg"`
}