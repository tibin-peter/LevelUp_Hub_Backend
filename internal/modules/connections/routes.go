package connections

import (
	"LevelUp_Hub_Backend/internal/middleware"
	"LevelUp_Hub_Backend/internal/modules/profile"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app fiber.Router, db *gorm.DB,jwtSecret string) {

	repo := NewRepository(db)
	mentorRepo:=profile.NewMentorRepository(db)
	service := NewService(repo,mentorRepo)
	handler := NewHandler(service)

	c := app.Group("/connections",middleware.AuthMiddleware(jwtSecret))

	c.Get("/my-mentors", handler.MyMentors)
	c.Get("/student-count",middleware.RequireRole("mentor"),handler.StudentCount)
	c.Get("/my-students",middleware.RequireRole("mentor"), handler.MyStudents)
}