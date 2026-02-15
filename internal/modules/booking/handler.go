package booking

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


//student side func

func (h *Handler) CreateBooking(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uint)
    var req CreateBookingRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.JSONError(c, 400, "invalid input")
    }

    // Capture the returned booking object
    booking, err := h.service.CreateBooking(userID, req.SlotID)
    if err != nil {
        return utils.JSONError(c, 400, err.Error())
    }

    // Return the booking ID so frontend can call POST /payments/order
    return utils.JSONSucess(c, "booking ready for payment", fiber.Map{
        "booking_id": booking.ID,
    })
}

func (h *Handler) GetStudentUpcoming(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	data, err := h.service.GetStudentUpcoming(userID)
	if err != nil {
		return  utils.JSONError(c,500,err.Error())
	}

	return utils.JSONSucess(c,"fetched",data)
}

func (h *Handler) GetStudentHistory(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	data, err := h.service.GetStudentHistory(userID)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "fetched", data)
}
func (h *Handler) CancelBooking(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return utils.JSONError(c, 400, "invalid booking id")
	}

	if err := h.service.CancelBooking(uint(id), userID); err != nil {
		return utils.JSONError(c, 400, err.Error())
	}

	return utils.JSONSucess(c, "booking cancelled", nil)
}


//mentor side functions

// approve booking
func (h *Handler) ApproveBooking(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.JSONError(c, 400, "invalid booking id")
	}

	if err := h.service.ApproveBooking(uint(id), userID); err != nil {
		return utils.JSONError(c, 400, err.Error())
	}

	return utils.JSONSucess(c, "booking approved", nil)
}

// reject booking
func (h *Handler) RejectBooking(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.JSONError(c, 400, "invalid booking id")
	}

	if err := h.service.RejectBooking(uint(id), userID); err != nil {
		return utils.JSONError(c, 400, err.Error())
	}

	return utils.JSONSucess(c, "booking rejected", nil)
}

// mentor upcoming
func (h *Handler) GetMentorUpcoming(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	data, err := h.service.GetMentorUpcoming(userID)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "fetched", data)
}

// mentor history
func (h *Handler) GetMentorHistory(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	data, err := h.service.GetMentorHistory(userID)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "fetched", data)
}