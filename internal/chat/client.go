package chat

import (
	"fmt"

	"github.com/gofiber/contrib/websocket"
)

type Message struct {
	IdOrigin      string `json:"id_origin"`
	IdDestination string `json:"id_destination"`
}

type ClientJson struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Client struct {
	Id            string
	Name          string
	Message       chan *Message
	Observer      *ChatObserver
	WebsocketConn *websocket.Conn
}

func NewClient(id string, name string, obs *ChatObserver, conn *websocket.Conn) *Client {
	return &Client{
		Id:            id,
		Name:          name,
		Observer:      obs,
		WebsocketConn: conn,
	}
}

func (c *Client) WriteMessage() {
	for {
		_, msg, _ := c.WebsocketConn.ReadMessage()
		fmt.Println(msg)
	}
}

func (c *Client) ReadMessage() {

	genericResponse := []byte("")

	for {
		c.WebsocketConn.WriteMessage(websocket.TextMessage, genericResponse)
	}
}
