package booking

import (
	"LevelUp_Hub_Backend/internal/middleware"
	"LevelUp_Hub_Backend/internal/modules/profile"
	"LevelUp_Hub_Backend/internal/modules/slot"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app fiber.Router, db *gorm.DB, jwtSecret string) {
	//Dependency wiring
	repo := NewRepository(db)
	slotrepo := slot.NewRepository(db)
	mentorrepo := profile.NewMentorRepository(db)
	service := NewService(repo, slotrepo, mentorrepo)
	handler := NewHandler(service)

	student := app.Group("/student", middleware.AuthMiddleware(jwtSecret), middleware.RequireRole("student"))

	student.Post("/bookings", handler.CreateBooking)
	student.Get("/bookings/upcoming", handler.GetStudentUpcoming)
	student.Get("/bookings/history", handler.GetStudentHistory)
	student.Patch("/bookings/:id/cancel", handler.CancelBooking)

	mentor := app.Group("/mentor", middleware.AuthMiddleware(jwtSecret), middleware.RequireRole("mentor"))

	mentor.Patch("/bookings/:id/approve", handler.ApproveBooking)
	mentor.Patch("/bookings/:id/reject", handler.RejectBooking)
	mentor.Get("/bookings/upcoming", handler.GetMentorUpcoming)
	mentor.Get("/bookings/history", handler.GetMentorHistory)
}
