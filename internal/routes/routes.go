package routes

import (
	"LevelUp_Hub_Backend/internal/modules/auth"
	"LevelUp_Hub_Backend/internal/modules/profile"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUp(
	app *fiber.App,
	db *gorm.DB,
	rdb *redis.Client,
	jwtSecret string,
	){
  
	api:=app.Group("/api")
	
	//module routes
	profile.RegisterRoutes(api,db,jwtSecret)
	auth.RegisterRoutes(api,db,rdb,jwtSecret)
}