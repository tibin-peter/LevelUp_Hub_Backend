package payment

import (
	"LevelUp_Hub_Backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	app fiber.Router,
	jwtSecret string,
	handler *Handler,
) {

	p := app.Group("/payments",
		middleware.AuthMiddleware(jwtSecret),
	)

	p.Post("/order", handler.CreateOrder)
	p.Post("/verify", handler.VerifyPayment)
	p.Get("/student", handler.GetStudentPayments)
	p.Get("/mentor/earnings", handler.GetMentorEarnings)
	p.Post("/withdraw", handler.RequestWithdraw)
	p.Get("/withdrawals", handler.GetMentorWithdrawals)

	// Admin routes
	admin := p.Group("/admin", middleware.RequireRole("admin"))
	admin.Get("/ledger", handler.AdminGetLedger)
	admin.Get("/payment-overview", handler.AdminGetPaymentOverview)
	admin.Get("/wallet-overview", handler.AdminGetWalletOverview)
	admin.Get("/wallet-transactions", handler.AdminGetWalletTransactions)
	admin.Get("/withdrawals", handler.AdminGetWithdrawals)
	admin.Patch("/withdraw/:id/approve", handler.AdminApproveWithdrawal)
	admin.Patch("/withdraw/:id/reject", handler.AdminRejectWithdrawal)
}