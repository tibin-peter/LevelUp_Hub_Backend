package ratings

import (
	"LevelUp_Hub_Backend/internal/middleware"
	"LevelUp_Hub_Backend/internal/modules/booking"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app fiber.Router, db *gorm.DB,jwtSecreat string) {
	//Dependency wiring
	repo := NewRepository(db)
	bookingrepo:=booking.NewRepository(db)
	service := NewService(repo,bookingrepo)
	handler := NewHandler(service)

	ratings := app.Group("/ratings",middleware.AuthMiddleware(jwtSecreat))

  ratings.Post("/", handler.CreateRating)
	ratings.Get("/mentor/:mentorId", handler.GetMentorRatings)
	ratings.Get("/mentor/:mentorId/summary", handler.GetMentorSummary)
	ratings.Get("/top", handler.GetTopMentors)
}