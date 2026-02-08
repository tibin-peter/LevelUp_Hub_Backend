package courses

import (
	"errors"
	"fmt"
)

type Service interface {
	ListCourses(filter CourseFilter) ([]Course,error)
	GetCourseByID(id uint) (*Course, error)
	ReplaceMentorCourses(mentorID uint,coursesIDs []uint)error
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

//replace mentor course
func(s *CourseService)ReplaceMentorCourses(mentorID uint,coursesIDs []uint)error{
	if len(coursesIDs)==0{
		return errors.New("at least one course needed")
	}

	//validate course exist
	for _,cid:=range coursesIDs{
		_,err:=s.repo.GetCourseByID(cid)
		if err!=nil{
			return fmt.Errorf("couses %d not found",cid)
		}
	}

	return s.repo.ReplaceMentorCourses(mentorID,coursesIDs)
}

//get course by mentor
func (s *CourseService) GetCoursesByMentor(mentorID uint,) ([]Course, error) {
	return s.repo.GetCoursesByMentor(mentorID)
}

//get mentor by course
func (s *CourseService) GetMentorsByCourse(courseID uint,) ([]MentorWithUser, error) {
	return s.repo.GetMentorsByCourse(courseID)
}