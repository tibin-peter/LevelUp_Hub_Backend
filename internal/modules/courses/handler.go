package courses

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

func (h *Handler) ListAllCourses(c *fiber.Ctx)error{
	var query CourseListQuery

	if err:=c.QueryParser(&query);err!=nil{
		return utils.JSONError(c,400,"invalid query")
	}

	filter := CourseFilter{
		Search:   query.Search,
		Category: query.Category,
		Level:    query.Level,
		Page:     query.Page,
		Limit:    query.Limit,
	}

	courses, err := h.service.ListCourses(filter)
	if err != nil {
		return utils.JSONError(c,500,err.Error())
	}

	return utils.JSONSucess(c,"course fetched",courses)
}

func (h *Handler) GetCourse(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.JSONError(c,400,"invalid id")
	}

	course, err := h.service.GetCourseByID(uint(id))
	if err != nil {
		return utils.JSONError(c,404,"course not found")
	}

	return utils.JSONSucess(c,"fetched successfully",course)
}

func (h *Handler) AddMentorCourse(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uint)

    var req AddMentorCourseRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.JSONError(c, 400, "Invalid data")
    }
		
    err := h.service.AddMentorCourse(userID, req.CourseID)
    if err != nil {
        return utils.JSONError(c, 400, err.Error())
    }

    return utils.JSONSucess(c, "Course added successfully", nil)
}

func (h *Handler) GetMentorCourses(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	courses, err := h.service.GetCoursesByMentor(userID)
	if err != nil {
		return utils.JSONError(c,500,err.Error())
	}

	return utils.JSONSucess(c,"success",courses)
}

func (h *Handler) GetMentorsByCourse(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.JSONError(c,400,"invalid id")
	}

	mentors, err := h.service.GetMentorsByCourse(uint(id))
	if err != nil {
		return utils.JSONError(c,500,err.Error())
	}

	return utils.JSONSucess(c,"success",mentors)
}