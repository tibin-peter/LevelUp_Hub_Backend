package payment

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

//create order

type CreateOrderRequest struct {
	BookingID uint `json:"booking_id"`
}

func (h *Handler) CreateOrder(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONError(c, 400, "invalid body")
	}

	res, err := h.service.CreateOrder(
		userID,
		req.BookingID,
	)
	if err != nil {
		return utils.JSONError(c, 400, err.Error())
	}

	return utils.JSONSucess(c, "order created", res)
}

//verify payment

func (h *Handler) VerifyPayment(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	var req VerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONError(c, 400, "invalid body")
	}

	if err := h.service.VerifyRequest(
		userID,
		req,
	); err != nil {
		return utils.JSONError(c, 400, err.Error())
	}

	return utils.JSONSucess(c, "payment verified", nil)
}


//student payments

func (h *Handler) GetStudentPayments(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	data, err := h.service.GetStudentPayment(userID)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "fetched", data)
}

//mentor earnings

func (h *Handler) GetMentorEarnings(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	data, err := h.service.GetMentorEarnings(userID)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}

	return utils.JSONSucess(c, "fetched", data)
}

//withdraw request

type WithdrawRequestDTO struct {
	Amount float64 `json:"amount"`
}

func (h *Handler) RequestWithdraw(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	var req WithdrawRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONError(c, 400, "invalid body")
	}

	if err := h.service.RequestWithdraw(
		userID,
		req.Amount,
	); err != nil {
		return utils.JSONError(c, 400, err.Error())
	}

	return utils.JSONSucess(c, "withdraw request submitted", nil)
}


// //admin list withdraws

// func (h *Handler) ListWithdrawRequests(c *fiber.Ctx) error {

// 	list, err := h.service.ListWithdrawRequests()
// 	if err != nil {
// 		return utils.JSONError(c, 500, err.Error())
// 	}

// 	return utils.JSONSucess(c, "fetched", list)
// }


// //admin approve withdraw

// func (h *Handler) ApproveWithdraw(c *fiber.Ctx) error {

// 	idParam := c.Params("id")
// 	id, _ := strconv.Atoi(idParam)

// 	if err := h.service.ApproveWithdraw(uint(id)); err != nil {
// 		return utils.JSONError(c, 400, err.Error())
// 	}

// 	return utils.JSONSucess(c, "withdraw approved", nil)
// }