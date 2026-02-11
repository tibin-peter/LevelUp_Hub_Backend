package auth

import (
	"LevelUp_Hub_Backend/internal/utils"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// sent otp func
func (h *Handler) SendOTP(c *fiber.Ctx) error {
	var sendotp SendOTPRequest
	if err := c.BodyParser(&sendotp); err != nil {
		return utils.JSONError(c, 400, "invalid request")
	}
	if sendotp.Email == "" {
		return utils.JSONError(c, 400, "email required")
	}
	if err := h.service.SendOTP(sendotp.Email); err != nil {
		return utils.JSONError(c, 500, err.Error())
	}
	return utils.JSONSucess(c, "otp sent",nil)
}

//func register
func (h *Handler) Register(c *fiber.Ctx) error {
	var required RegisterRequest
	if err := c.BodyParser(&required); err != nil {
		return utils.JSONError(c, 400, "invalid details")
	}
	log.Println("register data:",required)
	if required.Name == "" || required.Email == "" || required.Password == "" ||required.OTP==""{
		return utils.JSONError(c, 400, "all field is required")
	}
	if required.Role != "student" && required.Role != "mentor" {
		return utils.JSONError(c, 400, "invalid role")
	}
	access,refresh,userData,err := h.service.Register(required.Name, required.Email, required.Password, required.Role,required.OTP)
	if err!=nil{
		return utils.JSONError(c,400,err.Error())
	}
	//set cookies
	c.Cookie(&fiber.Cookie{
		Name: "access_token",
		Value: access,
		HTTPOnly: true,
		SameSite: "Lax",
		MaxAge: 7*24*3600,
	})

	c.Cookie(&fiber.Cookie{
		Name: "refresh_token",
		Value: refresh,
		HTTPOnly: true,
		SameSite: "Lax",
		MaxAge: 7*24*3600,
	})
	// Prepare the data for the frontend
    responseData := AuthResponseData{
        Email: userData.Email,
        Role:  userData.Role,
        Name:  userData.Name,
				IsVerified: userData.IsVerified,
				ProfilePicURL: userData.ProfilePicURL,
    }
	return utils.JSONSucess(c, "registed and logged in",responseData)
}

//func for login
func(h *Handler)Login(c *fiber.Ctx)error{
	var req LoginRequest
	if err:=c.BodyParser(&req);err!=nil{
		return utils.JSONError(c,400,"invalid details")
	}
	log.Println("dat:",req)
	if req.Email==""||req.Password==""{
		return utils.JSONError(c,400,"email and password required")
	}
	access,refresh,userData,err := h.service.Login(req.Email,req.Password)
	if err!=nil{
		return utils.JSONError(c,400,err.Error())
	}
	//set cookies
	c.Cookie(&fiber.Cookie{
		Name: "access_token",
		Value: access,
		HTTPOnly: true,
		SameSite: "Lax",
		MaxAge: 7*24*3600,
	})

	c.Cookie(&fiber.Cookie{
		Name: "refresh_token",
		Value: refresh,
		HTTPOnly: true,
		SameSite: "Lax",
		MaxAge: 7*24*3600,
	})
	// 4. Prepare Response Data
    responseData := AuthResponseData{
        Email: userData.Email,
        Role:  userData.Role,
        Name:  userData.Name,
				IsVerified: userData.IsVerified,
				ProfilePicURL: userData.ProfilePicURL,
    }
	return utils.JSONSucess(c, "loggin successfull",responseData)
}

//func for refresh 
func (h *Handler)Refresh(c *fiber.Ctx)error{
	refresh:=c.Cookies("refresh_token")
	if refresh==""{
		return utils.JSONError(c,401,"token is missing")
	}

	access,refresh,err := h.service.Refresh(refresh)
	if err!=nil{
		return utils.JSONError(c,400,err.Error())
	}
	//set cookies
	c.Cookie(&fiber.Cookie{
		Name: "access_token",
		Value: access,
		HTTPOnly: true,
		SameSite: "Lax",
		MaxAge: 7*24*3600,
	})

	c.Cookie(&fiber.Cookie{
		Name: "refresh_token",
		Value: refresh,
		HTTPOnly: true,
		SameSite: "Lax",
		MaxAge: 7*24*3600,
	})
	return utils.JSONSucess(c,"token refreshed",nil)
}

//func for logout
func(h *Handler)Logout(c *fiber.Ctx)error{
	userID:=c.Locals("userID").(uint)

	if err:=h.service.Logout(userID);err!=nil{
		return utils.JSONError(c,400,err.Error())
	}
	c.Cookie(&fiber.Cookie{
		Name: "access_token",
		Value: "",
		HTTPOnly: true,
		Expires: time.Now().Add(-time.Hour),
	})
	c.Cookie(&fiber.Cookie{
		Name: "refresh_token",
		Value: "",
		HTTPOnly: true,
		Expires: time.Now().Add(-time.Hour),
	})
	return utils.JSONSucess(c,"logout successfull",nil)
}