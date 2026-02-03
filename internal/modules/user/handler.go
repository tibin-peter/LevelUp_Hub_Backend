package user

import (
	"LevelUp_Hub_Backend/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// get by id
func (h *Handler) GetUserById(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return utils.JSONError(c, 400, "invalid user id")
	}

	user, err := h.service.GetUserById(uint(id))
	if err != nil {
		return utils.JSONError(c, 404, err.Error())
	}
	return utils.JSONSucess(c, user)
}

//Get by email
func (h *Handler) GetUserByEmail(c *fiber.Ctx) error {
	email := c.Params("email")

	if email == "" {
		return utils.JSONError(c, 400, "invalid user email")
	}
	user, err := h.service.FindByEmail(email)
	if err != nil {
		return utils.JSONError(c, 404, err.Error())
	}
	return utils.JSONSucess(c, user)
}

//Update user
func (h *Handler) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return utils.JSONError(c, 400, "invalid user id")
	}
	var dto UpdateUserDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.JSONError(c, 400, "invalid required body")
	}
	user, err := h.service.GetUserById(uint(id))
	if err != nil {
		return utils.JSONError(c, 404, err.Error())
	}
	if dto.Name != "" {
		user.Name = dto.Name
	}
	if dto.ProfilePicURL != "" {
		user.ProfilePicURL = dto.ProfilePicURL
	}
	if err := h.service.Update(user); err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "user update")
}

//Delete User
func (h *Handler)DeleteUser(c *fiber.Ctx)error{
	idParam:=c.Params("id")
	id,err:=strconv.Atoi(idParam)
	if err!=nil||id<=0{
		return utils.JSONError(c,400,"invalid id")
	}
	err1:=h.service.Delete(uint(id))
	if err1!=nil{
		return utils.JSONError(c,404,err1.Error())
	}
	return utils.JSONSucess(c,"user deleted successfully")
}