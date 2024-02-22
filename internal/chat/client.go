package chat

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/contrib/websocket"
)

type Message struct {
	IdOrigin      string `json:"id_origin"`
	IdDestination string `json:"id_destination"`
	Content       string `json:"content"`
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
		Id:                 id,
		Name:               name,
		Observer:           obs,
		WebsocketConn:      conn,
		ReceiveMessageChan: make(chan string), // TOO IMPORTANT (If there isn't an channel initialized, the message will never be received)
	}
}

func (c *Client) ReadMessageFromClient() {

	defer func() {
		// c.Observer.Unregister <- c
		c.Observer.UnsubscribeClientChan <- c
		_ = c.WebsocketConn.Close()
	}()

	for {
		_, msg, err := c.WebsocketConn.ReadMessage()

		if err != nil {
			fmt.Println(err)
			break
		}

		chatMessage := Message{}
		json.Unmarshal(msg, &chatMessage)
		fmt.Println("MESSAGE RECEIVED!")
		fmt.Println(chatMessage.Content)
		fmt.Println(chatMessage.IdDestination)
		fmt.Println("---------------------")
		chatMessage.IdOrigin = c.Id
		c.Observer.SendMessageChan <- &chatMessage
	}
}

func (c *Client) WriteMessageToClient() {

	fmt.Println("Goroutine write message to client starts")
	// c.WebsocketConn.WriteMessage(websocket.TextMessage, []byte("Este es un mensaje de prueba"))

	defer func() {
		_ = c.WebsocketConn.Close()
	}()

	for {
		select {
		case messageReceived := <-c.ReceiveMessageChan:
			fmt.Println("SENDING TO CLIENT")
			fmt.Println(messageReceived)
			c.WebsocketConn.WriteMessage(websocket.TextMessage, []byte(messageReceived))
		}
	}

	// for messageReceived := range c.ReceiveMessageChan {
	// 	c.WebsocketConn.WriteMessage(websocket.TextMessage, []byte(messageReceived))
	// }
}
