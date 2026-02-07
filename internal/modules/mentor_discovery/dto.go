package mentordiscovery

//for giving respose for the mentor details
type MentorCard struct {
	MentorProfileID uint    `json:"mentor_profile_id"`
	UserID          uint    `json:"user_id"`
	Name            string  `json:"name"`
	ProfilePicURL   string  `json:"profile_Pic_url"`
	Bio             string  `json:"bio"`
	Skills          string  `json:"skills"`
	HourlyPrice     float64 `json:"hourly_price"`
	RatingAvg       float64 `json:"rating_avg"`
	ExperienceYears int     `json:"experience_years"`
}

//dto for giving the filter details
type MentorFilter struct {
    Skill     string  `query:"skill"`
    Search    string  `query:"search"`
    MinPrice  float64 `query:"min_price"`
    MaxPrice  float64 `query:"max_price"`
    MinRating float64 `query:"min_rating"`
    Sort      string  `query:"sort"`
    Page      int     `query:"page"`
    Limit     int     `query:"limit"`
    Offset    int
}