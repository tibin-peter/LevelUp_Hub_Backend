package auth

import (
	"LevelUp_Hub_Backend/internal/middleware"
	"LevelUp_Hub_Backend/internal/modules/profile"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterRoutes(
	r fiber.Router,
  db *gorm.DB,
	rdb *redis.Client,
	jwtSecret string,
	){

	//dependency wiring
	repo:=profile.NewRepository(db)
	mentorrepo:=profile.NewMentorRepository(db)
	service:=NewService(repo,mentorrepo,rdb,jwtSecret)
	handler:=NewHandler(service)

	//group routes
	auth:=r.Group("/auth")

	auth.Post("/send-otp",handler.SendOTP)
	auth.Post("/register",handler.Register)
	auth.Post("/login",handler.Login)
	auth.Post("/refresh",handler.Refresh)

	//protected
	auth.Post("/logout",middleware.AuthMiddleware(jwtSecret),handler.Logout)


}