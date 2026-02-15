package profile

import "time"

type UpdateUserDTO struct {
	Name          string `json:"name"`
	ProfilePicURL string
}

type MentorProfileInput struct {
	Bio             string  `json:"bio"`
	Category        string  `json:"category"`
	SubCategory     string  `json:"sub_category"`
	Skills          string  `json:"skills"`
	Languages       string  `json:"languages"`
	HourlyPrice     float64 `json:"hourly_price"`
	ExperienceYears int     `json:"experience_years"`
	LinkedinUrl     string  `json:"linkedin_url"`
  Documents       string  `json:"documents"`
}

// pulic mentor response
type MentorProfileResponse struct {
	ID     uint
	UserID uint
	Name   string
	ProfilePicURL string

	Bio       string
	Skills    string
	Languages string

	HourlyPrice     float64
	ExperienceYears int

	RatingAvg    float64
	TotalReviews int

	CreatedAt time.Time
	UpdatedAt time.Time
}