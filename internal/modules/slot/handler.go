package slot

import (
	"LevelUp_Hub_Backend/internal/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// create slot
func (h *Handler) CreateSlot(c *fiber.Ctx)error{
   var req CreateSlotRequest
	 if err:=c.BodyParser(&req);err!=nil{
		return utils.JSONError(c,400,"invalid request")
	 }
	 mentorID:=c.Locals("userID").(uint)
	 err:=h.service.CreateSlot(mentorID,req)
	 if err!=nil{
		return utils.JSONError(c,400,err.Error())
	 }
	 return utils.JSONSucess(c,"slot created",nil)
}

//get mentor slot
func (h *Handler)GetMentorSlots(c *fiber.Ctx)error{
	mentorID,ok:=c.Locals("userID").(uint)
	if !ok{
		return utils.JSONError(c,401,"unauthorized")
	}
	slots,err:=h.service.GetMentorSlots(mentorID)
	if err!=nil{
		return utils.JSONError(c,500,err.Error())
	}
	return utils.JSONSucess(c,"slots fetched",slots)
}

//delete slot
func (h *Handler)DeleteSlot(c *fiber.Ctx)error{
	mentorID,ok:=c.Locals("userID").(uint)
	if !ok{
		return utils.JSONError(c,401,"unauthorized")
	}
	slotIDStr:=c.Params("slot_id")
	id,_:=strconv.Atoi(slotIDStr)
	slotID:=uint(id)
	if err:=h.service.DeleteSlot(mentorID,slotID);err!=nil{
		return utils.JSONError(c,500,err.Error())
	}
	return utils.JSONSucess(c,"slot deleted",nil)
}

//get available slot
func (h *Handler)GetAvailableSlots(c *fiber.Ctx)error{
	mentorParam:=c.Params("mentor_id")
	mentorInt,err:=strconv.Atoi(mentorParam)
	if err!=nil{
		return utils.JSONError(c,400,"invalid mentor id")
	}
	mentorID:=uint(mentorInt)

	//for date if select
	dateStr:=c.Query("date")
	var parseDate *time.Time
	if dateStr!=""{
		d,err:=time.Parse("2006-01-02",dateStr)
		if err!=nil{
			return utils.JSONError(c,400,"choose a relevent date")
		}
		parseDate = &d
	}
	slots, err := h.service.GetAvailableSlots(mentorID, parseDate)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "available slots fetched", slots)
}