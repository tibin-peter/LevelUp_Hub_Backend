package dashboard

import (
	"LevelUp_Hub_Backend/internal/middleware"
	"LevelUp_Hub_Backend/internal/modules/booking"
	"LevelUp_Hub_Backend/internal/modules/connections"
	"LevelUp_Hub_Backend/internal/modules/favorites"
	"LevelUp_Hub_Backend/internal/modules/payment"
	"LevelUp_Hub_Backend/internal/modules/profile"
	"LevelUp_Hub_Backend/internal/modules/ratings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app fiber.Router, db *gorm.DB, jwtSecret string) {
	//Dependency wiring
	b:=booking.NewRepository(db)
	fav:=favorites.NewRepository(db)
	r:=ratings.NewRepository(db)
	pay:=payment.NewRepository(db)
	c:=connections.NewRepository(db)
	m:=profile.NewMentorRepository(db)
	service := NewService(b,fav,r,pay,c,m)
	handler := NewHandler(service)

	dashboard := app.Group("/dashboard",middleware.AuthMiddleware(jwtSecret))

	// student dashboard
	dashboard.Get("/student",middleware.RequireRole("student"),handler.StudentDashboard)

	// student dashboard
	dashboard.Get("/mentor",middleware.RequireRole("mentor"),handler.MentorDashboard)
}