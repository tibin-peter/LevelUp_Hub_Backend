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

	// // admin later if needed not now
	// admin := p.Group("/admin", auth,
	// 	middleware.RequireRole("admin"),
	// )

	// admin.Get("/withdraws", handler.ListWithdrawRequests)
	// admin.Patch("/withdraw/:id", handler.ApproveWithdraw)

}