package usermanagement

import (
	"LevelUp_Hub_Backend/internal/middleware"
	"LevelUp_Hub_Backend/internal/modules/complaints"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app fiber.Router, db *gorm.DB, jwtSecret string) {
	//Dependency wiring
	repo := NewRepository(db)
	complaintsRepo:=complaints.NewRepository(db)
	service := NewService(repo,complaintsRepo)
	handler := NewHandler(service)

	admin := app.Group(
		"/admin",
		middleware.AuthMiddleware(jwtSecret),
		middleware.RequireRole("admin"),
	)

	// users
	admin.Get("/students",handler.Students)
	admin.Get("/mentors",handler.Mentors)
	admin.Patch("/users/:id/block",handler.BlockUser)
	admin.Patch("/users/:id/unblock",handler.UnblockUser)

	// mentor approvals
	admin.Get("/mentor-approvals",handler.PendingMentors)
	admin.Patch("/mentor-approvals/:id/approve",handler.ApproveMentor)
	admin.Patch("/mentor-approvals/:id/reject",handler.RejectMentor)

	// complaints
	admin.Get("/complaints",handler.Complaints)
	admin.Put("/complaints/:id/reply",handler.ReplyComplaint)

}