package favorites

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

func (h *Handler) AddFavorite(c *fiber.Ctx) error {

	user := c.Locals("userID")
	studentID, ok := user.(uint)
	if !ok {
		return utils.JSONError(c, 401, "unauthorized")
	}

	id, _ := strconv.Atoi(c.Params("mentorProfileId"))

	err := h.service.AddFavorite(studentID, uint(id))
	if err != nil {
		return utils.JSONError(c, 400, err.Error())
	}

	return utils.JSONSucess(c, "added to favorites", nil)
}


func (h *Handler) RemoveFavorite(c *fiber.Ctx) error {

	user := c.Locals("userID")
	studentID := user.(uint)

	id, _ := strconv.Atoi(c.Params("mentorProfileId"))

	err := h.service.RemoveFavorite(studentID, uint(id))
	if err != nil {
		return utils.JSONError(c, 400, err.Error())
	}

	return utils.JSONSucess(c, "removed from favorites", nil)
}

func (h *Handler) ListFavorites(c *fiber.Ctx) error {

	user := c.Locals("userID")
	studentID := user.(uint)

	res, err := h.service.ListFavorites(studentID)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "my favorites", res)
}