package middleware

import (
	"LevelUp_Hub_Backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

//for jwt authentication
func AuthMiddleware(JWTSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Try to read from cookie
		tokenStr := c.Cookies("access_token")

		// 2. Try to read from Authorization header if cookie is missing
		if tokenStr == "" {
			authHeader := c.Get("Authorization")
			if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenStr = authHeader[7:]
			}
		}

		if tokenStr == "" {
			return utils.JSONError(c, 401, "missing token")
		}

		// validate token
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
