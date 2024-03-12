package handlers

import (
	"fmt"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pelusa-v/pelusa-chat.git/internal/chat"
)

func RegisterRoomViewHandler(c *fiber.Ctx) error {

	if c.Method() == fiber.MethodPost {
		nickName := c.FormValue("nick")
		return c.Redirect(fmt.Sprintf("/room/%s", nickName))
	}

	return c.Render("internal/views/register.html", nil)
}

func ChatRoomViewHandler(c *fiber.Ctx) error {
	data2 := struct {
		Items []string
	}{
		Items: []string{"Apple", "Banana", "Orange"},
	}

	data := fiber.Map{
		"nick":    c.Params("nick"),
		"clients": data2,
	}
	return c.Render("internal/views/room.html", data)
}

func RegisterHandler(c *websocket.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)

	client := chat.NewClient(uuid.New().String(), c.Params("nick"), &chat.Manager, c)
	client.Manager.SubscribeClientChan <- client

	var registerNotification = &chat.WebSocketData{
		Notification: chat.RegisteringNotification{
			ClientId:     client.Id,
			ClientName:   client.Name,
			Registerting: true,
			// Content:    fmt.Sprintf("%s has enter to this chat room", client.Name),
		},
		IsMessage: false,
	}
	client.Manager.BroadcastNotificationChan <- registerNotification

	go client.ReadMessageFromClient()
	go client.WriteMessageToClient()

	wg.Wait()
}

func ShowClientsHandler(c *fiber.Ctx) error {
	var clients []chat.ClientJson
	for _, client := range chat.Manager.Clients {
		clients = append(clients, chat.ClientJson{Id: client.Id, Name: client.Name})
	}

	return c.JSON(clients)
}
