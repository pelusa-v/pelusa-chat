package chat

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/contrib/websocket"
)

type WebSocketData struct {
	Message      Message
	Notification ClientNotification
}

type Message struct {
	IdOrigin      string `json:"id_origin"`
	IdDestination string `json:"id_destination"`
	Content       string `json:"content"`
}

type NotificationType string

const (
	RegisterNotification   NotificationType = "register"
	UnregisterNotification NotificationType = "unregister"
)

type ClientNotification struct {
	ClientId   string           `json:"client_id"`
	ClientName string           `json:"client_name"`
	Content    string           `json:"content"`
	Type       NotificationType `json:"type"`
}

type ClientJson struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Client struct {
	Id                 string
	Name               string
	Manager            *ChatManager
	WebsocketConn      *websocket.Conn
	ReceiveMessageChan chan *WebSocketData
}

func NewClient(id string, name string, manager *ChatManager, conn *websocket.Conn) *Client {
	return &Client{
		Id:                 id,
		Name:               name,
		Manager:            manager,
		WebsocketConn:      conn,
		ReceiveMessageChan: make(chan *WebSocketData), // TOO IMPORTANT (If there isn't an channel initialized, the message will never be received)
	}
}

func (c *Client) ReadMessageFromClient() {

	defer func() {
		// c.Observer.Unregister <- c
		c.Manager.UnsubscribeClientChan <- c
		_ = c.WebsocketConn.Close()
	}()

	for {
		_, msg, err := c.WebsocketConn.ReadMessage()

		if err != nil {
			fmt.Println(err)
			break
		}

		chatMessage := WebSocketData{}
		json.Unmarshal(msg, &chatMessage)
		fmt.Println("MESSAGE RECEIVED!")
		fmt.Println(chatMessage.Message.Content)
		fmt.Println(chatMessage.Message.IdDestination)
		fmt.Println("---------------------")
		chatMessage.Message.IdOrigin = c.Id
		c.Manager.SendMessageChan <- &chatMessage
	}
}

func (c *Client) WriteMessageToClient() {

	fmt.Println("Goroutine write message to client starts")

	defer func() {
		_ = c.WebsocketConn.Close()
	}()

	for {
		select {
		case messageReceived := <-c.ReceiveMessageChan:
			fmt.Println("SENDING TO CLIENT")
			fmt.Println(messageReceived)
			data, _ := json.Marshal(messageReceived)
			c.WebsocketConn.WriteMessage(websocket.TextMessage, data)
		}
	}
}
