package adminprofile

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

func (h *Handler) GetAdminProfile(c *fiber.Ctx) error{
	id:=c.Locals("userID").(uint)
	profile,err:=h.service.GetAdminProfile(id)
	if err!=nil{
		return utils.JSONError(c,400,err.Error())
	}
	data:=&AdminProfileResponse{
		Name: profile.Name,
		Email: profile.Email,
		Role: profile.Role,
		ProfilePicURL: profile.ProfilePicURL,
	}
	return utils.JSONSucess(c,"success",data)
}

func (h *Handler)UpdateProfile(c *fiber.Ctx)error{
	id:=c.Locals("userID").(uint)
	var req UpdateProfile
	if err:=c.BodyParser(&req);err!=nil{
		return utils.JSONError(c,400,"invalid details")
	}
	profile,err:=h.service.UpdateProfile(id,req)
	if err!=nil{
		return utils.JSONError(c,500,err.Error())
	}
	data:=&AdminProfileResponse{
		Name: profile.Name,
		Email: profile.Email,
		Role: profile.Role,
		ProfilePicURL: profile.ProfilePicURL,
	}
	return utils.JSONSucess(c,"updated successfully",data)
}

func (h *Handler)UpdateProfilePicture(c *fiber.Ctx)error{
	id:=c.Locals("userID").(uint)
	var req UpdateProfilePicture
	if err:=c.BodyParser(&req);err!=nil{
		return utils.JSONError(c,400,"invalid details")
	}
	profile,err:=h.service.UpdateProfilePicture(id,req)
	if err!=nil{
		return utils.JSONError(c,500,err.Error())
	}
	data:=&AdminProfilePicResponse{
		ProfilePicURL: profile.ProfilePicURL,
	}
	return utils.JSONSucess(c,"updated successfully",data)
}

func (h *Handler)ChangePassword(c *fiber.Ctx)error{
	id:=c.Locals("userID").(uint)
	var req ChangePassword
	if err:=c.BodyParser(&req);err!=nil{
		return utils.JSONError(c,400,"invalid details")
	}
	err:=h.service.ChangePassword(id,req)
	if err!=nil{
		return utils.JSONError(c,500,"failed to update")
	}
	return utils.JSONSucess(c,"password updated successfully",nil)
}