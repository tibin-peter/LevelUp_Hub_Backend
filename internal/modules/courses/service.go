package courses

import (
	"fmt"
)

type Service interface {
	ListCourses(filter CourseFilter) ([]Course,error)
	GetCourseByID(id uint) (*Course, error)
	AddMentorCourse(mentorID uint, courseID uint) error
	GetCoursesByMentor(mentorID uint,) ([]Course, error)
	GetMentorsByCourse(courseID uint,) ([]MentorWithUser, error)
}

type CourseService struct {
	repo Repository
}

func NewService(r Repository) *CourseService {
	return &CourseService{repo: r}
}

// list course
func (s *CourseService) ListCourses(filter CourseFilter) ([]Course,error){
	//default pagination for safety
	if filter.Page<=0{
		filter.Page=1
	}

	if filter.Limit<=0||filter.Limit>50{
		filter.Limit=10
	}

	return s.repo.ListAllCourses(filter)
}

//get course by id
func (s *CourseService) GetCourseByID(id uint) (*Course, error) {
	return s.repo.GetCourseByID(id)
}

func (s *CourseService) AddMentorCourse(mentorID uint, courseID uint) error {
    //  Validate course exists
    _, err := s.repo.GetCourseByID(courseID)
    if err != nil {
        return fmt.Errorf("course %d not found", courseID)
    }

    return s.repo.AddMentorCourse(mentorID, courseID)
}

//get course by mentor
func (s *CourseService) GetCoursesByMentor(mentorID uint,) ([]Course, error) {
	return s.repo.GetCoursesByMentor(mentorID)
}

//get mentor by course
func (s *CourseService) GetMentorsByCourse(courseID uint,) ([]MentorWithUser, error) {
	return s.repo.GetMentorsByCourse(courseID)
}