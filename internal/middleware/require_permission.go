package middleware

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RequirePermission(targetSlug string, db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Get user role from JWT/Session (Assuming it's stored in locals)
		roleName, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// 2. Query the junction table to see if this role has the slug
		var count int64
		db.Table("role_permissions").
			Joins("JOIN roles ON roles.id = role_permissions.role_id").
			Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
			Where("roles.name = ? AND permissions.slug = ?", roleName, targetSlug).
			Count(&count)

		if count == 0 {
			return c.Status(403).JSON(fiber.Map{
				"error": "Access Denied: You do not have the permission: " + targetSlug,
			})
		}

		return c.Next()
	}
}