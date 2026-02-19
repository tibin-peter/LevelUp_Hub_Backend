package admindashboard


type AdminDashboardDTO struct {
	TotalStudents int64   `json:"total_students"`
	TotalMentors  int64   `json:"total_mentors"`
	ActiveCourses int64   `json:"active_courses"`
	TotalRevenue  float64 `json:"total_revenue"`

	RevenueChart []RevenueChartDTO `json:"revenue_chart"`
}

type RevenueChartDTO struct {
	Label   string  `json:"label"`
	Revenue float64 `json:"revenue"`
}