package ratings

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

func (h *Handler) CreateRating(c *fiber.Ctx)error{
	studentID := c.Locals("userID").(uint)

	var req CreateRatingRequest
	if err:=c.BodyParser(&req);err!=nil{
		return  utils.JSONError(c,400,"invalid body")
	}
	res,err:=h.service.CreateRating(studentID,req)
	if err!=nil{
     return utils.JSONError(c,400,err.Error())
	}	

	return utils.JSONSucess(c,"Rating created successfully",res)
}

func (h *Handler) GetMentorRatings(c *fiber.Ctx) error {

	mentorID, err := strconv.Atoi(c.Params("mentorId"))
	if err != nil {
		return utils.JSONError(c, 400, "invalid mentor id")
	}

	res, err := h.service.GetMentorRatings(uint(mentorID))
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "mentor ratings", res)
}

func (h *Handler) GetMentorSummary(c *fiber.Ctx) error {

	mentorID, err := strconv.Atoi(c.Params("mentorId"))
	if err != nil {
		return utils.JSONError(c, 400, "invalid mentor id")
	}

	res, err := h.service.GetMentorSummary(uint(mentorID))
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "mentor summary", res)
}

func (h *Handler) GetTopMentors(c *fiber.Ctx) error {

	res, err := h.service.GetTopMentors()
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "top mentors", res)
}