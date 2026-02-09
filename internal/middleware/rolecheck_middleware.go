package middleware

import (
	"LevelUp_Hub_Backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func RequireRole(reqRole string) fiber.Handler{
	return func (c *fiber.Ctx)error {
		role,ok:=c.Locals("role").(string)
		if !ok{
			return utils.JSONError(c,401,"unauthorized")
		}
		if role!=reqRole{
			return utils.JSONError(c,401,"unauthorized")
		}
		return c.Next()
	}
}