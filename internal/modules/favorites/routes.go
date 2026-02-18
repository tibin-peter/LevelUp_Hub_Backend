package favorites

import (
	"LevelUp_Hub_Backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app fiber.Router, db *gorm.DB, jwtSecreat string) {
	// Dependency wiring
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	fav := app.Group("/favorite",middleware.AuthMiddleware(jwtSecreat))

	fav.Post("/:mentorProfileId", handler.AddFavorite)
	fav.Delete("/:mentorProfileId", handler.RemoveFavorite)
	fav.Get("/my", handler.ListFavorites)
}