package utils

import "github.com/gofiber/fiber/v2"


// Define a standard response struct
type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"` // omitempty hides the field if data is nil
}

func JSONSucess(c *fiber.Ctx, message string, data interface{}) error {
    return c.JSON(Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

//func for error response
func JSONError(c *fiber.Ctx,status int ,message string)error{
	return c.Status(status).JSON(fiber.Map{
		"success":false,
		"error":message,
	})
}