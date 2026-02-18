package connections

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

func (h *Handler) MyMentors(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	res, err := h.service.GetMyMentors(userID)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "connected mentors", res)
}

func (h *Handler) StudentCount(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	count, err := h.service.GetStudentCount(userID)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "connected students",
		CountResponse{Total: count})
}

func (h *Handler) MyStudents(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	res, err := h.service.GetMyStudents(userID)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "connected students", res)
}