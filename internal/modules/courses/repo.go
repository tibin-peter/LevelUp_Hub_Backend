package courses

import (
	"LevelUp_Hub_Backend/internal/repository/generic"

	"gorm.io/gorm"
)

type Repository interface {
	ListAllCourses(filter CourseFilter) ([]Course, error)
	GetCourseByID(courseID uint) (*Course, error)
	AddMentorCourse(mentorID uint, courseID uint) error
	GetCoursesByMentor(mentorID uint) ([]Course, error)
  GetMentorsByCourse(courseID uint) ([]MentorWithUser, error)
}

type repo struct {
	base *generic.Repository[Course]
	db   *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{
		base: generic.NewRepository[Course](db),
		db:   db,
	}
}

//struct for mentor course and user related info
type MentorWithUser struct {   
	UserID        uint
	Name          string
	ProfilePicURL string
	ID            uint
	Bio           string
	HourlyPrice   float64
	RatingAvg     float64
}

//func for list all couses with filter
func (r *repo) ListAllCourses(filter CourseFilter) ([]Course, error) {
	var courses []Course

	//bind the stuct to the table
	query := r.db.Model(&Course{})

	//search by tittle
	if filter.Search != "" {
		query = query.Where("LOWER(tittle) LIKE LOWER(?)", "%"+filter.Search+"%")
	}

	//by category
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}

	//by level
	if filter.Level != "" {
		query = query.Where("level = ?", filter.Level)
	}

	//panination
	limit := 10
	if filter.Limit > 0 {
		limit = filter.Limit
	}

	offset := 0
	if filter.Page > 0 {
		offset = (filter.Page - 1) * limit
	}

	query = query.Limit(limit).Offset(offset)

	err := query.Find(&courses).Error

	return courses, err

}

//func for get course by id
func (r *repo) GetCourseByID(courseID uint) (*Course, error) {
	var course Course
	err := r.db.First(&course, courseID).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *repo) AddMentorCourse(mentorID uint, courseID uint) error {
    mapping := MentorCourse{
        MentorID: mentorID,
        CourseID: courseID,
    }

    // .FirstOrCreate prevents duplicate entries and SQL errors 
    // if the user clicks the same course twice.
    err := r.db.Where(MentorCourse{MentorID: mentorID, CourseID: courseID}).
                FirstOrCreate(&mapping).Error
                
    return err
}

//func for get courses by mentor
func (r *repo) GetCoursesByMentor(mentorID uint) ([]Course, error) {
	var courses []Course

	err := r.db.Model(&Course{}).
		Joins("JOIN mentor_courses mc ON mc.course_id = courses.id").
		Where("mc.mentor_id = ?", mentorID).
		Find(&courses).Error

	return courses, err
}

//fucn for get mentor by course
func (r *repo) GetMentorsByCourse(courseID uint) ([]MentorWithUser, error) {
	var mentors []MentorWithUser

	err := r.db.Table("mentor_courses mc").
		Select(`
						u.id as user_id,
						u.name,
						u.profile_Pic_url,
						mp.id as id,
						mp.bio,
						mp.hourly_price,
						mp.rating_avg
						`).
		Joins("JOIN users u ON u.id = mc.mentor_id").
		Joins("JOIN mentor_profiles mp ON mp.user_id = u.id").
		Where("mc.course_id = ?", courseID).
		Scan(&mentors).Error

	return mentors, err
}
