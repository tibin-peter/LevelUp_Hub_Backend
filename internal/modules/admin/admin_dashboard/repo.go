package admindashboard

import (
	"LevelUp_Hub_Backend/internal/modules/booking"
	"LevelUp_Hub_Backend/internal/modules/courses"
	"LevelUp_Hub_Backend/internal/modules/payment"
	"LevelUp_Hub_Backend/internal/modules/profile"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	CountStudents() (int64, error)
	CountMentors() (int64, error)
	CountOfCourses() (int64, error)
	CountEnrollments() (int64, error)
	TotalPlatformRevenue(start, end time.Time) (float64, error)
	RevenueChart(start, end time.Time, filter string) ([]RevenueChartDTO, error)
	GetRecentActivities() ([]RecentActivityDTO, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

// total students count
func (r *repo) CountStudents() (int64, error) {
	var count int64

	err := r.db.
		Model(&profile.User{}).
		Where("role = ?", "student").
		Count(&count).Error

	return count, err
}

// total mentors count
func (r *repo) CountMentors() (int64, error) {
	var count int64

	err := r.db.
		Model(&profile.User{}).
		Where("role = ?", "mentor").
		Count(&count).Error

	return count, err
}

//total course count
func (r *repo) CountOfCourses() (int64, error) {
	var count int64

	err := r.db.
		Model(&courses.Course{}).
		Count(&count).Error

	return count, err
}

func (r *repo) CountEnrollments() (int64, error) {
	var count int64
	err := r.db.Model(&booking.Booking{}).Where("status IN ?", []string{"paid", "confirmed", "completed"}).Count(&count).Error
	return count, err
}

func (r *repo) TotalPlatformRevenue(start, end time.Time) (float64, error) {

	var total float64

	err := r.db.
		Model(&payment.Payment{}).
		Where("status = ?", "paid").
		Where("created_at BETWEEN ? AND ?", start, end).
		Select("COALESCE(SUM(amount),0)").
		Scan(&total).Error

	return total, err
}

//for revenue chart
func (r *repo) RevenueChart(start, end time.Time, filter string) ([]RevenueChartDTO, error) {

	var result []RevenueChartDTO

	var query string

	switch filter {

	case "week":
		query = `
		SELECT
		  TO_CHAR(created_at,'Dy') as label,
		  COALESCE(SUM(amount),0) as revenue
		FROM payments
		WHERE status='paid'
		AND created_at BETWEEN ? AND ?
		GROUP BY label
		ORDER BY MIN(created_at)
		`

	case "month":
		query = `
		SELECT
		  TO_CHAR(created_at,'DD Mon') as label,
		  COALESCE(SUM(amount),0) as revenue
		FROM payments
		WHERE status='paid'
		AND created_at BETWEEN ? AND ?
		GROUP BY label
		ORDER BY MIN(created_at)
		`

	case "year":
		query = `
		SELECT
		  TO_CHAR(created_at,'Mon') as label,
		  COALESCE(SUM(amount),0) as revenue
		FROM payments
		WHERE status='paid'
		AND created_at BETWEEN ? AND ?
		GROUP BY label
		ORDER BY MIN(created_at)
		`
	}

	err := r.db.Raw(query, start, end).
		Scan(&result).Error

	return result, err
}

func (r *repo) GetRecentActivities() ([]RecentActivityDTO, error) {
	var activities []RecentActivityDTO

	// 1. Get Newest Users
	var newUsers []profile.User
	r.db.Order("created_at desc").Limit(5).Find(&newUsers)
	for _, u := range newUsers {
		activities = append(activities, RecentActivityDTO{
			Type:        "user",
			Title:       "New User Entry",
			Description: fmt.Sprintf("%s joined as %s", u.Name, u.Role),
			Time:        u.CreatedAt.Format("15:04 PM"),
		})
	}

	// 2. Get Newest Courses
	var newCourses []courses.Course
	r.db.Order("created_at desc").Limit(3).Find(&newCourses)
	for _, c := range newCourses {
		activities = append(activities, RecentActivityDTO{
			Type:        "course",
			Title:       "New Curriculum Node",
			Description: fmt.Sprintf("Course '%s' initialized", c.Title),
			Time:        c.CreatedAt.Format("15:04 PM"),
		})
	}

	return activities, nil
}