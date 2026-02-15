package booking

import (
	"LevelUp_Hub_Backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
    app fiber.Router,
    jwtSecret string,
    handler *Handler,
) {

    student := app.Group("/student",
        middleware.AuthMiddleware(jwtSecret),
        middleware.RequireRole("student"),
    )

    student.Post("/bookings", handler.CreateBooking)
    student.Get("/bookings/upcoming", handler.GetStudentUpcoming)
    student.Get("/bookings/history", handler.GetStudentHistory)
    student.Patch("/bookings/:id/cancel", handler.CancelBooking)

    mentor := app.Group("/mentor",
        middleware.AuthMiddleware(jwtSecret),
        middleware.RequireRole("mentor"),
    )

    mentor.Patch("/bookings/:id/approve", handler.ApproveBooking)
    mentor.Patch("/bookings/:id/reject", handler.RejectBooking)
    mentor.Get("/bookings/upcoming", handler.GetMentorUpcoming)
    mentor.Get("/bookings/history", handler.GetMentorHistory)
}