package payment

import (
	"LevelUp_Hub_Backend/internal/utils"
	"log"
	"strconv"

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
	log.Println("payment created")

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

func (h *Handler) GetMentorWithdrawals(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	data, err := h.service.GetMentorWithdrawals(userID)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}
	return utils.JSONSucess(c, "fetched withdrawals", data)
}

// Admin handlers

func (h *Handler) AdminGetLedger(c *fiber.Ctx) error {
	search := c.Query("search", "")
	status := c.Query("status", "")
	data, err := h.service.GetAdminPaymentLedger(search, status)
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}
	return utils.JSONSucess(c, "fetched ledger", data)
}

func (h *Handler) AdminGetPaymentOverview(c *fiber.Ctx) error {
	data, err := h.service.GetAdminPaymentOverview()
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}
	return utils.JSONSucess(c, "fetched summary", data)
}

func (h *Handler) AdminGetWalletOverview(c *fiber.Ctx) error {
	data, err := h.service.GetAdminWalletOverview()
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}
	return utils.JSONSucess(c, "fetched wallet summary", data)
}

func (h *Handler) AdminGetWalletTransactions(c *fiber.Ctx) error {
	data, err := h.service.GetAdminWalletTransactions()
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}
	return utils.JSONSucess(c, "fetched transactions", data)
}

func (h *Handler) AdminGetWithdrawals(c *fiber.Ctx) error {
	data, err := h.service.GetAdminWithdrawals()
	if err != nil {
		return utils.JSONError(c, 500, err.Error())
	}
	return utils.JSONSucess(c, "fetched withdrawals", data)
}

func (h *Handler) AdminApproveWithdrawal(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.ApproveWithdrawal(uint(id)); err != nil {
		return utils.JSONError(c, 400, err.Error())
	}
	return utils.JSONSucess(c, "withdrawal approved", nil)
}

func (h *Handler) AdminRejectWithdrawal(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.RejectWithdrawal(uint(id)); err != nil {
		return utils.JSONError(c, 400, err.Error())
	}
	return utils.JSONSucess(c, "withdrawal rejected", nil)
}