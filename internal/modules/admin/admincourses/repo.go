package admincourses

import (
	"LevelUp_Hub_Backend/internal/modules/courses"

	"gorm.io/gorm"
)

type Repository interface {
	CountCourses() (int64, error)

	CreateCourse(course *courses.Course) error

	UpdateCourse(id uint, input courses.Course) error

	DeleteCourse(id uint) error

	ListAllCourses(filter CourseFilter) ([]courses.Course, int64, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

func (r *repo) CountCourses() (int64, error) {
	var count int64
	err := r.db.Model(&courses.Course{}).Count(&count).Error
	return count, err
}

func (r *repo) CreateCourse(course *courses.Course) error {
	err := r.db.Create(course).Error
	return err
}
func (r *repo) UpdateCourse(id uint, input courses.Course) error {
	err := r.db.
		Model(&courses.Course{}).
		Where("id = ?", id).
		Updates(input).Error
	return err
}

func (r *repo) DeleteCourse(id uint) error {
	err := r.db.Where("id = ?", id).Delete(&courses.Course{}).Error
	return err
}

func (r *repo) ListAllCourses(filter CourseFilter) ([]courses.Course, int64, error) {

	var course []courses.Course
	var total int64

	query := r.db.Model(&courses.Course{})

	//  Search
	if filter.Search != "" {
		query = query.Where(
			"title ILIKE ?",
			"%"+filter.Search+"%",
		)
	}

	//  Category filter
	if filter.Category != "" {
		query = query.Where(
			"category = ?",
			filter.Category,
		)
	}

	//  level filter
	if filter.Level != "" {
		query = query.Where(
			"level = ?",
			filter.Level,
		)
	}

	// count before pagination
	query.Count(&total)

	// pagination
	limit := 10
	if filter.Limit > 0 {
		limit = filter.Limit
	}

	page := 1
	if filter.Page > 0 {
		page = filter.Page
	}

	offset := (page - 1) * limit

	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&course).Error

	return course, total, err
}
