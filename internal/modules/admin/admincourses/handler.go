package admincourses

import (
	"LevelUp_Hub_Backend/internal/modules/courses"
	"LevelUp_Hub_Backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

//count

func (h *Handler) CountCourses(c *fiber.Ctx) error {

	count, err := h.service.CountCourses()
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "course count",count)
}

// list courses (Search + Filter + Pagination)

func (h *Handler) ListCourses(c *fiber.Ctx) error {

	filter := CourseFilter{}

	if err := c.QueryParser(&filter); err != nil {
		return utils.JSONError(c, 400, "invalid query")
	}

	data, total, err := h.service.ListCourses(filter)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "courses", fiber.Map{
		"items": data,
		"total": total,
		"page":  filter.Page,
		"limit": filter.Limit,
	})
}

// create course

func (h *Handler) CreateCourse(c *fiber.Ctx) error {

	var body courses.Course

	if err := c.BodyParser(&body); err != nil {
		return utils.JSONError(c, 400, "invalid body")
	}

	err := h.service.CreateCourse(body)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "course created", nil)
}

//update

func (h *Handler) UpdateCourse(c *fiber.Ctx) error {

	id, _ := c.ParamsInt("id")

	var body courses.Course

	if err := c.BodyParser(&body); err != nil {
		return utils.JSONError(c, 400, "invalid body")
	}

	err := h.service.UpdateCourse(uint(id), body)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "course updated", nil)
}

//delete

func (h *Handler) DeleteCourse(c *fiber.Ctx) error {

	id, _ := c.ParamsInt("id")

	err := h.service.DeleteCourse(uint(id))
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "course deleted", nil)
}