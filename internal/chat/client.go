package chat

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/contrib/websocket"
)

type Message struct {
	IdOrigin        string `json:"id_origin"`
	IdDestination   string `json:"id_destination"`
	NameOrigin      string `json:"name_origin"`
	NameDestination string `json:"name_destination"`
	Content         string `json:"content"`
	Broadcast       bool   `json:"broadcast"`
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
	ReceiveMessageChan chan *Message
}

func NewClient(id string, name string, manager *ChatManager, conn *websocket.Conn) *Client {
	return &Client{
		Id:                 id,
		Name:               name,
		Manager:            manager,
		WebsocketConn:      conn,
		ReceiveMessageChan: make(chan *Message), // TOO IMPORTANT (If there isn't an channel initialized, the message will never be received)
	}
}

func (c *Client) ReadMessageFromClient() {

	defer func() {
		c.Manager.UnsubscribeClientChan <- c
		_ = c.WebsocketConn.Close()

		var unregisterNotification = &Message{
			IdOrigin:   "Manager",
			NameOrigin: "Manager",
			Content:    fmt.Sprintf("***  %s + ( + %s + ) left this room ***", c.Name, c.Id),
			Broadcast:  true,
		}

		c.Manager.BroadcastNotificationChan <- unregisterNotification
	}()

	for {
		_, msg, err := c.WebsocketConn.ReadMessage()

		if err != nil {
			fmt.Println(err)
			break
		}

		chatMessage := Message{}
		json.Unmarshal(msg, &chatMessage)
		chatMessage.IdOrigin = c.Id
		chatMessage.Broadcast = false
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
