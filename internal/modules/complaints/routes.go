package complaints

import (
	"LevelUp_Hub_Backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app fiber.Router, db *gorm.DB, jwtSecret string) {

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	complaints := app.Group("/complaints", middleware.AuthMiddleware(jwtSecret))

	// USER
	complaints.Post("/", handler.CreateComplaint)
	complaints.Get("/my", handler.MyComplaints)

}