package middleware

import (
	"LevelUp_Hub_Backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(JWTSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		//read cookie
		tokenStr := c.Cookies("access_token")
		if tokenStr == "" {
			return utils.JSONError(c, 401, "missing token")
		}
		//validate token
		claims, err := utils.ValidateToken(tokenStr, JWTSecret)
		if err != nil {
			return utils.JSONError(c, 401, "invalid token")
		}
		//store in ctx
		c.Locals("userID", claims.UserID)
		c.Locals("role", claims.Role)
		c.Locals("email", claims.Email)
		return c.Next()
	}
}
