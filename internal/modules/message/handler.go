package message

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Handler struct{
	service Service
	hub *Hub
}

func NewHandler(s Service,h *Hub)*Handler{
	return &Handler{s,h}
}

func (h *Handler) HandleWS(c *websocket.Conn,userID uint){

	cid64,_ := strconv.ParseUint(
		c.Params("conversationID"),
		10,32,
	)
	cid := uint(cid64)

	ok,_ := h.service.IsUserInConversation(cid,userID)
	if !ok { 
        fmt.Println(" access Denied for user", userID)
        c.Close()
        return 
    }
	if !ok { c.Close(); return }

	client := &Client{
		UserID:userID,
		Conn:c,
		Send:make(chan interface{}),
	}

	h.hub.register<-client
	defer func(){
		h.hub.unregister<-client
		c.Close()
	}()

	// writer
	go func(){
		for msg:=range client.Send{
			c.WriteJSON(msg)
		}
	}()

	// reader
	for{
		var payload SendWSMessage

		if err:=c.ReadJSON(&payload);err!=nil{
			break
		}
		if payload.Content==""{continue}

		h.service.SendMessage(cid,userID,payload.Content)

		other,_ := h.service.GetOtherUser(cid,userID)

		data := fiber.Map{
			"conversation_id":cid,
			"from":userID,
			"content":payload.Content,
		}

		h.hub.forward<-HubMessage{
			ToUserID:other,
			Data:data,
		}

		client.Send<-fiber.Map{"status":"sent"}
	}
}

func (h *Handler) StartConversation(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	var req StartConversationRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(400, "invalid body")
	}

	id, err := h.service.StartConversation(
		userID,
		req.TargetUserID,
	)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{
		"conversation_id": id,
	})
}

func (h *Handler) ListConversations(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	data, err := h.service.ListUserConversations(userID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(data)
}

func (h *Handler) GetMessages(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	cid, _ := strconv.Atoi(c.Params("id"))
	before, _ := strconv.Atoi(c.Query("before", "0"))

	ok, _ := h.service.IsUserInConversation(uint(cid), userID)
	if !ok {
		return fiber.ErrForbidden
	}

	msgs, err := h.service.GetMessages(
		uint(cid),
		uint(before),
	)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(msgs)
}