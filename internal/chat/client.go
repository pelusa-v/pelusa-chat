package chat

import (
	"fmt"

	"github.com/gofiber/contrib/websocket"
)

type Message struct {
	IdOrigin      string `json:"id_origin"`
	IdDestination string `json:"id_destination"`
	Content       string
}

type ClientJson struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Client struct {
	Id                 string
	Name               string
	Observer           *ChatObserver
	WebsocketConn      *websocket.Conn
	ReceiveMessageChan chan string
}

func NewClient(id string, name string, obs *ChatObserver, conn *websocket.Conn) *Client {
	return &Client{
		Id:            id,
		Name:          name,
		Observer:      obs,
		WebsocketConn: conn,
	}
}

func (c *Client) ReadMessageFromClient() {
	for {
		_, msg, _ := c.WebsocketConn.ReadMessage()
		fmt.Println(msg)
	}
}

func (c *Client) WriteMessageToClient() {

	defer c.WebsocketConn.Close()

	for messageReceived := range c.ReceiveMessageChan {
		c.WebsocketConn.WriteMessage(websocket.TextMessage, []byte(messageReceived))
	}
}
