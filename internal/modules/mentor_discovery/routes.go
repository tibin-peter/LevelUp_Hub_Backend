package mentordiscovery

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app fiber.Router,db *gorm.DB){
	//Dependency wiring
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	//routes
	app.Get("/mentors",handler.GetMentors)
}