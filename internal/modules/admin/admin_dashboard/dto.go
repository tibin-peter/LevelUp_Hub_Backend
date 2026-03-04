package admindashboard


type AdminDashboardDTO struct {
	TotalStudents int64   `json:"total_students"`
	TotalMentors  int64   `json:"total_mentors"`
	ActiveCourses int64   `json:"active_courses"`
	TotalRevenue  float64 `json:"total_revenue"`
	TotalEnrollments int64  `json:"total_enrollments"`

	RevenueChart []RevenueChartDTO `json:"revenue_chart"`
	RecentActivities []RecentActivityDTO `json:"recent_activities"`
}

type RecentActivityDTO struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Time        string `json:"time"`
}

type RevenueChartDTO struct {
	Label   string  `json:"label"`
	Revenue float64 `json:"revenue"`
}