package courses

//struct for filter
type CourseFilter struct {
	Search   string `query:"search"`
	Category string `query:"category"`
	Level    string `query:"level"`
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
}

//for query
type CourseListQuery struct {
	Search   string `query:"search"`
	Category string `query:"category"`
	Level    string `query:"level"`
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
}

//update course details
type AddMentorCourseRequest struct {
	CourseID uint `json:"course_id" validate:"required"`
}

//dto for response
type CourseResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Level       string `json:"level"`
}

//dto for mentor details
type MentorSummaryResponse struct {
	UserID        uint    `json:"user_id"`
	Name          string  `json:"name"`
	ProfilePicURL string  `json:"profile_pic_url"`
	Bio           string  `json:"bio"`
	HourlyPrice   float64 `json:"hourly_price"`
	RatingAvg     float64 `json:"rating_avg"`
}

//dto for mentorcourse res
type MentorCourseResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}
