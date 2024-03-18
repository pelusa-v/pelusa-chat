package chat

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/contrib/websocket"
)

type Message struct {
	OriginId        string `json:"origin_id"`
	DestinationId   string `json:"destination_id"`
	OriginName      string `json:"origin_name"`
	DestinationName string `json:"destination_name"`
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
	Manager            *ChatManager    // Used to pass data through Manager channels or to perform Manager actions
	WebsocketConn      *websocket.Conn // websocket connection used by client to communicate with server
	ReceiveMessageChan chan *Message   // channel through which messages are received
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

func (c *Client) WriteMessages() {

	defer func() {
		c.Manager.UnsubscribeClientChan <- c
		_ = c.WebsocketConn.Close()

		var unregisterNotification = &Message{
			OriginId:   "Manager",
			OriginName: "Manager",
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
		chatMessage.OriginId = c.Id
		fmt.Println(string(msg))
		c.Manager.SendMessageChan <- &chatMessage
	}
}

func (c *Client) ReadMessages() {

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
