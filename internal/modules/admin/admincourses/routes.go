package admincourses

import (
	"LevelUp_Hub_Backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app fiber.Router, db *gorm.DB, jwtSecret string) {
	//Dependency wiring
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	admin := app.Group(
		"/admin",
		middleware.AuthMiddleware(jwtSecret),
		middleware.RequireRole("admin"),
	)

	admin.Get("/courses/count", handler.CountCourses)

	admin.Get("/courses", handler.ListCourses)

	admin.Post("/courses", handler.CreateCourse)

	admin.Put("/courses/:id", handler.UpdateCourse)

	admin.Delete("/courses/:id", handler.DeleteCourse)

}
