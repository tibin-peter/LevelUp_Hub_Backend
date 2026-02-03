package user

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App,db *gorm.DB){
	//Dependency wiring
	repo:=NewRepository(db)
	service:=NewService(repo)
	handler:=NewHandler(service)

	//Route group
	userGrop:=app.Group("/users")

	userGrop.Get("/:id",handler.GetUserById)
	userGrop.Get("/email/:email",handler.GetUserByEmail)
	userGrop.Get("/update/:id",handler.UpdateUser)
	userGrop.Get("/delete/:id",handler.DeleteUser)
}