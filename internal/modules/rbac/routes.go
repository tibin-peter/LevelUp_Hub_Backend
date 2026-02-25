package rbac

import (
	"LevelUp_Hub_Backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, jwtSecret string) {
	//Dependency wiring
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	permission := r.Group("/admin", middleware.AuthMiddleware(jwtSecret), middleware.RequireRole("admin"))

	permission.Post("/permission", handler.CreatePermission)
	permission.Get("/permissions", handler.ListPermissions)
	permission.Patch("/roles/toggle", handler.ToggleRolePermission)
	permission.Post("/roles", handler.CreateRole)
	permission.Get("/roles", handler.ListRoles)
	permission.Get("/roles/:role/permissions", handler.GetPermissionsByRole)

}
