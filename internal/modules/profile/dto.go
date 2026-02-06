package profile


type UpdateUserDTO struct{
	Name     string `json:"name"`
	ProfilePicURL string
}

type MentorProfileInput struct {
    Bio             string  `json:"bio"`
    Skills          string  `json:"skills"`
    Languages       string  `json:"languages"`
    HourlyPrice     float64 `json:"hourly_price"`
    ExperienceYears int     `json:"experience_years"`
}