package complaints

import (
	"LevelUp_Hub_Backend/internal/utils"
	// "strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateComplaint(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)
	role := c.Locals("role").(string)

	var req CreateComplaintRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONError(c, 400, "invalid body")
	}

	err := h.service.CreateComplaint(userID, role, req)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "complaint submitted", nil)
}

func (h *Handler) MyComplaints(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	list, err := h.service.MyComplaints(userID)
	if err != nil {
		return utils.JSONError(c,500,err.Error())
	}

	return utils.JSONSucess(c,"my complaints",list)
}