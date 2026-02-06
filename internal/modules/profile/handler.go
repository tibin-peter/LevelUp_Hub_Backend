package profile

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
	return utils.JSONSucess(c,"successfull", user)
}

//Get by email
func (h *Handler) GetUserByEmail(c *fiber.Ctx) error {
	email := c.Params("email")

	if email == "" {
		return utils.JSONError(c, 400, "invalid user email")
	}
	user, err := h.service.FindUserByEmail(email)
	if err != nil {
		return utils.JSONError(c, 404, err.Error())
	}
	return utils.JSONSucess(c,"successfull", user)
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
	if err := h.service.UpdateUser(user); err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "user update",user)
}

//Delete User
func (h *Handler)DeleteUser(c *fiber.Ctx)error{
	idParam:=c.Params("id")
	id,err:=strconv.Atoi(idParam)
	if err!=nil||id<=0{
		return utils.JSONError(c,400,"invalid id")
	}
	err1:=h.service.DeleteUser(uint(id))
	if err1!=nil{
		return utils.JSONError(c,404,err1.Error())
	}
	return utils.JSONSucess(c,"user deleted successfully",nil)
}

//create mentor profile
func (h *Handler) CreateMentorProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input MentorProfileInput
	if err := c.BodyParser(&input); err != nil {
		return utils.JSONError(c,400,"invalid input")
	}

	profile,err := h.service.CrateMentorProfile(userID, input)
	if err != nil {
		return utils.JSONError(c,400,err.Error())
	}
  
	return utils.JSONSucess(c,"mentor profile created",profile)
}

//get mentor profile
func (h *Handler) GetMentorProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	profile, err := h.service.GetMentorProfileByUserID(userID)
	if err != nil {
		return utils.JSONError(c,404,"profile not found")
	}

	return utils.JSONSucess(c,"fetched successfully",profile)
}

//update mentor profile
func (h *Handler) UpdateMentorProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input MentorProfileInput
	if err := c.BodyParser(&input); err != nil {
		return utils.JSONError(c,400,"invalid input")
	}

	profile,err := h.service.UpdateMentorProfile(userID, input)
	if err != nil {
		return utils.JSONError(c,400,err.Error())
	}

	return utils.JSONSucess(c,"profile updated",profile)
}

//get public mentor profile
func (h *Handler) GetPublicMentorProfile(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.JSONError(c,400,"invalid id")
	}

	profile, err := h.service.GetPublicMentorProfile(uint(id))
	if err != nil {
		return utils.JSONError(c,404,"mentor not found")
	}

	return utils.JSONSucess(c,"fetched successfully",profile)
}