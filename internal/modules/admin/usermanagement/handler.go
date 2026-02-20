package usermanagement

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

func (h *Handler) Students(c *fiber.Ctx) error {

	search := c.Query("search", "")

	data, err := h.service.ListUsers("student", search)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "students", data)
}

func (h *Handler) Mentors(c *fiber.Ctx) error {

	search := c.Query("search","")

	data,err := h.service.ListUsers("mentor",search)
	if err != nil {
		return utils.JSONError(c,500,err.Error())
	}

	return utils.JSONSucess(c,"mentors",data)
}

func (h *Handler) BlockUser(c *fiber.Ctx) error {

	id,_ := c.ParamsInt("id")

	err := h.service.BlockUser(uint(id))
	if err != nil {
		return utils.JSONError(c,500,err.Error())
	}

	return utils.JSONSucess(c,"user blocked",nil)
}

func (h *Handler) UnblockUser(c *fiber.Ctx) error {

	id,_ := c.ParamsInt("id")

	err := h.service.UnblockUser(uint(id))
	if err != nil {
		return utils.JSONError(c,500,err.Error())
	}

	return utils.JSONSucess(c,"user unblocked",nil)
}

func (h *Handler) PendingMentors(c *fiber.Ctx) error {

	data,err := h.service.PendingMentors()
	if err != nil {
		return utils.JSONError(c,500,err.Error())
	}

	return utils.JSONSucess(c,"pending mentors",data)
}

func (h *Handler) ApproveMentor(c *fiber.Ctx) error {

	id,_ := c.ParamsInt("id")

	err := h.service.ApproveMentor(uint(id))
	if err != nil {
		return utils.JSONError(c,500,err.Error())
	}

	return utils.JSONSucess(c,"mentor approved",nil)
}

func (h *Handler) RejectMentor(c *fiber.Ctx) error {

	id,_ := c.ParamsInt("id")

	err := h.service.RejectMentor(uint(id))
	if err != nil {
		return utils.JSONError(c,500,err.Error())
	}

	return utils.JSONSucess(c,"mentor rejected",nil)
}


func (h *Handler) Complaints(c *fiber.Ctx) error {

	search := c.Query("search","")

	data,err := h.service.ListComplaints(search)
	if err != nil {
		return utils.JSONError(c,500,err.Error())
	}

	return utils.JSONSucess(c,"complaints",data)
}

func (h *Handler) ReplyComplaint(c *fiber.Ctx) error {

	id,_ := c.ParamsInt("id")

	var body struct{
		Reply string `json:"reply"`
		Status string `json:"status"`
	}

	if err := c.BodyParser(&body); err != nil {
		return utils.JSONError(c,400,"invalid body")
	}

	err := h.service.ReplyComplaint(uint(id),body.Reply,body.Status)
	if err != nil {
		return utils.JSONError(c,500,err.Error())
	}

	return utils.JSONSucess(c,"reply added",nil)
}