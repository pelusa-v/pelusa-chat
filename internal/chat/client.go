package chat

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/contrib/websocket"
)

type WebSocketData struct {
	Message      Message
	Notification RegisteringNotification
	IsMessage    bool
}

type Message struct {
	IdOrigin      string `json:"id_origin"`
	IdDestination string `json:"id_destination"`
	Content       string `json:"content"`
}

type RegisteringNotification struct {
	ClientId     string `json:"client_id"`
	ClientName   string `json:"client_name"`
	Registerting bool   `json:"registering"`
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
		c.Manager.UnsubscribeClientChan <- c
		_ = c.WebsocketConn.Close()

		var uregisterNotification = &WebSocketData{
			Notification: RegisteringNotification{
				ClientId:     c.Id,
				ClientName:   c.Name,
				Registerting: false,
			},
			IsMessage: false,
		}
		c.Manager.BroadcastNotificationChan <- uregisterNotification
	}()

	for {
		_, msg, err := c.WebsocketConn.ReadMessage()

		if err != nil {
			fmt.Println(err)
			break
		}

		chatMessage := WebSocketData{}
		json.Unmarshal(msg, &chatMessage)
		chatMessage.Message.IdOrigin = c.Id
		fmt.Println(string(msg))
		c.Manager.SendMessageChan <- &chatMessage
	}
}

func (c *Client) WriteMessageToClient() {

	defer func() {
		_ = c.WebsocketConn.Close()
	}()

	for {
		select {
		case messageReceived := <-c.ReceiveMessageChan:
			data, _ := json.Marshal(messageReceived)
			c.WebsocketConn.WriteMessage(websocket.TextMessage, data)
		}
	}
}
