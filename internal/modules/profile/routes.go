package profile

import (
	"LevelUp_Hub_Backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB,jwtSecret string) {
	//Dependency wiring
	repo := NewRepository(db)
	mentorrepo := NewMentorRepository(db)
	service := NewService(repo, mentorrepo)
	handler := NewHandler(service)

	//Route group user
	userGrop := r.Group("/users")

	userGrop.Get("/:id", handler.GetUserById)
	userGrop.Get("/email/:email", handler.GetUserByEmail)
	userGrop.Put("/update/:id", handler.UpdateUser)
	userGrop.Delete("/delete/:id", handler.DeleteUser)

	//Route group mentor
	mentor := r.Group("/mentor", middleware.AuthMiddleware(jwtSecret))

	mentor.Post("/profile", handler.CreateMentorProfile)
	mentor.Get("/profile", handler.GetMentorProfile)
	mentor.Put("/profile", handler.UpdateMentorProfile)

	r.Get("/mentors/:id", handler.GetPublicMentorProfile)
}
