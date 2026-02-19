package admindashboard

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

func (h *Handler) Dashboard(c *fiber.Ctx) error {

	filter := c.Query("filter", "month")

	data, err := h.service.Dashboard(filter)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "admin dashboard", data)
}