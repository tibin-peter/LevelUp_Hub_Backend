package courses

import (
	"LevelUp_Hub_Backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app fiber.Router, db *gorm.DB,jwtSecret string) {
	//Dependency wiring
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	course := app.Group("/courses")

	course.Get("/", handler.ListAllCourses)
	course.Get("/:id", handler.GetCourse)
	course.Get("/:id/mentors", handler.GetMentorsByCourse)

	mentor := app.Group("/mentor", middleware.AuthMiddleware(jwtSecret))

	mentor.Post("/courses", handler.AddMentorCourse)
	mentor.Get("/courses", handler.GetMentorCourses)

}
