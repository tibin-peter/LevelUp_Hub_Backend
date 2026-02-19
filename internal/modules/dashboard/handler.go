package dashboard

import (
	"LevelUp_Hub_Backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) StudentDashboard(c *fiber.Ctx) error {

	userVal := c.Locals("userID")
	userID, ok := userVal.(uint)
	if !ok {
		return utils.JSONError(c, 401, "unauthorized")
	}

	data, err := h.service.StudentDashboard(userID)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "student dashboard", data)
}

func (h *Handler) MentorDashboard(c *fiber.Ctx) error {

	userVal := c.Locals("userID")
	userID, ok := userVal.(uint)
	if !ok {
		return utils.JSONError(c, 401, "unauthorized")
	}

	data, err := h.service.MentorDashboard(userID)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "mentor dashboard", data)
}