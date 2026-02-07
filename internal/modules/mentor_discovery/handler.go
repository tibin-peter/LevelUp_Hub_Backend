package mentordiscovery

import (
	"LevelUp_Hub_Backend/internal/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// func for get mentors
func (h *Handler) GetMentors(c *fiber.Ctx)error{

	filter:=new(MentorFilter)
	if err:=c.QueryParser(filter);err!=nil{
		return utils.JSONError(c,400,"Invalid query params")
	}
	fmt.Printf("%+v\n", filter)

	mentors,err:=h.service.GetMentors(filter)
	if err!=nil{
		return utils.JSONError(c,500,"failed to fetch mentors")
	}
	return utils.JSONSucess(c,"fetched successfully",mentors)
}