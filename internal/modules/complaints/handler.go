package complaints

import (
	"LevelUp_Hub_Backend/internal/utils"
	"strconv"

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

func (h *Handler) AllComplaints(c *fiber.Ctx) error {

	list, err := h.service.AllComplaints()
	if err != nil {
		return utils.JSONError(c,500,err.Error())
	}

	return utils.JSONSucess(c,"all complaints",list)
}

func (h *Handler) AdminReply(c *fiber.Ctx) error {

	id,_ := strconv.Atoi(c.Params("id"))

	var req ReplyRequest
	c.BodyParser(&req)

	err := h.service.AdminReply(
		uint(id),
		req.Reply,
		req.Status,
	)

	if err != nil {
		return utils.JSONError(c,500,err.Error())
	}

	return utils.JSONSucess(c,"replied",nil)
}