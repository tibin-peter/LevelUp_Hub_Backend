package message

import (
	"LevelUp_Hub_Backend/internal/middleware"
	"LevelUp_Hub_Backend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(
	app fiber.Router,
	db *gorm.DB,
	jwtSecret string,
){

	repo := NewRepository(db)
	service := NewService(repo)

	hub := NewHub()
	go hub.Run()

	handler := NewHandler(service, hub)

	msg := app.Group("/messages")

	// http
msg.Post("/conversation",
	middleware.AuthMiddleware(jwtSecret),
	handler.StartConversation,
)

msg.Get("/listconversations",
	middleware.AuthMiddleware(jwtSecret),
	handler.ListConversations,
)

msg.Get("/:id/messages",
	middleware.AuthMiddleware(jwtSecret),
	handler.GetMessages,
)
  
	//websocket
	msg.Get(
		"/ws/:conversationID",
		websocket.New(func(c *websocket.Conn) {

			token := c.Query("token")

			if token == ""{
				token = c.Cookies("access_token")
			}

			claims, err := utils.ValidateToken(token, jwtSecret)
			if err != nil {
				c.Close()
				return
			}

			handler.HandleWS(c,claims.UserID)
		}),
	)
}