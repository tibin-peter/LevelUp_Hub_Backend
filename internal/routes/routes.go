package routes

import (
	"LevelUp_Hub_Backend/internal/modules/user"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetUp(app *fiber.App,db *gorm.DB){
	user.RegisterRoutes(app,db)
}