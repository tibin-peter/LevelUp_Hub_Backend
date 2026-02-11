package message

import (
	"log"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	UserID uint
	Conn   *websocket.Conn
	Send   chan interface{}
}

type HubMessage struct {
	ToUserID uint
	Data     interface{}
}

type Hub struct {
	clients    map[uint]*Client
	register   chan *Client
	unregister chan *Client
	forward    chan HubMessage
}

func NewHub()*Hub{
	return &Hub{
		clients:make(map[uint]*Client),
		register:make(chan *Client),
		unregister:make(chan *Client),
		forward:make(chan HubMessage),
	}
}

func (h *Hub) Run(){
	for{
		select{

		case c:=<-h.register:
			h.clients[c.UserID]=c
			log.Println("register",c.UserID)

		case c:=<-h.unregister:
			if _,ok:=h.clients[c.UserID];ok{
				delete(h.clients,c.UserID)
				close(c.Send)
			}
			log.Println("unregister",c.UserID)

		case msg:=<-h.forward:
			if client,ok:=h.clients[msg.ToUserID];ok{
				client.Send<-msg.Data
			}
			log.Println("send",msg.ToUserID)
		}
	}
}