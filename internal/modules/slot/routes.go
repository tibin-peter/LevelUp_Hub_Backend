package slot

import (
	"LevelUp_Hub_Backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app fiber.Router,db *gorm.DB,jwtSecret string){
	//dependecy wiring
	repo:=NewRepository(db)
	service:=NewService(repo)
	handler:=NewHandler(service)

	//student routes
	student:=app.Group("/mentors",middleware.AuthMiddleware(jwtSecret),middleware.RequireRole("student"))
	student.Get("/:mentor_id/slots",handler.GetAvailableSlots)

	//mentor routes
	mentor:=app.Group("/mentor",middleware.AuthMiddleware(jwtSecret),middleware.RequireRole("mentor"))
	mentor.Post("/slots",handler.CreateSlot)
	mentor.Get("/slots",handler.GetMentorSlots)
	mentor.Delete("/slots/:slot_id",handler.DeleteSlot)
}
