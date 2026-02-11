package routes

import (
	"LevelUp_Hub_Backend/internal/modules/auth"
	"LevelUp_Hub_Backend/internal/modules/booking"
	"LevelUp_Hub_Backend/internal/modules/courses"
	"LevelUp_Hub_Backend/internal/modules/mentor_discovery"
	"LevelUp_Hub_Backend/internal/modules/message"
	"LevelUp_Hub_Backend/internal/modules/profile"
	"LevelUp_Hub_Backend/internal/modules/slot"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUp(
	app *fiber.App,
	db *gorm.DB,
	rdb *redis.Client,
	jwtSecret string,
	){


	// upgrade middleware globali
	app.Use("/api/messages/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
  
	api:=app.Group("/api")
	
	//module routes
	profile.RegisterRoutes(api,db,jwtSecret)
	auth.RegisterRoutes(api,db,rdb,jwtSecret)
	mentordiscovery.RegisterRoutes(api,db)
	courses.RegisterRoutes(api,db,jwtSecret)
	slot.RegisterRoutes(api,db,jwtSecret)
	booking.RegisterRoutes(api,db,jwtSecret)
	message.RegisterRoutes(api,db,jwtSecret)
}