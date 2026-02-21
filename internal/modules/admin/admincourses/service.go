package admincourses

import (
	"LevelUp_Hub_Backend/internal/modules/courses"
	"fmt"
)

type Service interface {
	CountCourses()(int64,error)

	CreateCourse(input courses.Course) error

	UpdateCourse(id uint,input courses.Course) error

	DeleteCourse(id uint) error

	ListCourses(filter CourseFilter)([]courses.Course,int64,error)
}

type service struct {
	repo       Repository
}

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) CountCourses()(int64,error){
	return s.repo.CountCourses()
}

func (s *service) CreateCourse(input courses.Course) error {

	if input.Title == "" {
		return fmt.Errorf("title required")
	}

	return s.repo.CreateCourse(&input)
}

func (s *service) UpdateCourse(id uint,input courses.Course) error {

	return s.repo.UpdateCourse(id,input)
}

func (s *service) DeleteCourse(id uint) error {
	return s.repo.DeleteCourse(id)
}

func (s *service) ListCourses(filter CourseFilter)([]courses.Course,int64,error){

	return s.repo.ListAllCourses(filter)
}