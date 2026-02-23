package adminprofile

import (
	"LevelUp_Hub_Backend/internal/middleware"
	"LevelUp_Hub_Backend/internal/modules/profile"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app fiber.Router, db *gorm.DB, jwtSecret string) {
	//Dependency wiring
	userRepo:=profile.NewRepository(db)
	service := NewService(userRepo)
	handler := NewHandler(service)

	admin := app.Group("/admin",middleware.AuthMiddleware(jwtSecret),middleware.RequireRole("admin"))
	admin.Get("/profile",handler.GetAdminProfile)
	admin.Put("/profile/update",handler.UpdateProfile)
	admin.Put("/profile/updateimage",handler.UpdateProfilePicture)
	admin.Put("/profile/password",handler.ChangePassword)
}