package utils

import "github.com/gofiber/fiber/v2"

//func for success response
func JSONSucess(c *fiber.Ctx,data interface{})error{
	return c.JSON(fiber.Map{
		"success":true,
		"data":data,
	})
}

//func for error response
func JSONError(c *fiber.Ctx,status int ,message string)error{
	return c.Status(status).JSON(fiber.Map{
		"success":false,
		"error":message,
	})
}