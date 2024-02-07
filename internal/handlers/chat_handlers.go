package handlers

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pelusa-v/pelusa-chat.git/internal/chat"
)

func RegisterRoomHandler(c *fiber.Ctx) error {

	if c.Method() == fiber.MethodPost {
		nickName := c.FormValue("nick")
		return c.Redirect(fmt.Sprintf("/room/%s", nickName))
	}

	return c.Render("internal/views/register.html", nil)
}

func ChatRoomHandler(c *fiber.Ctx) error {
	data := fiber.Map{
		"nick": c.Params("nick"),
	}
	return c.Render("internal/views/room.html", data)
}

func RegisterHandler(c *websocket.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)

	obs := c.Locals("observer").(*chat.ChatObserver)
	client := chat.NewClient(uuid.New().String(), c.Params("nick"), obs, c)
	client.Observer.SubscribeClientChan <- client

	go client.ReadMessageFromClient()
	go client.WriteMessageToClient()

	wg.Wait()
}

func BroacastHandler(c *websocket.Conn) {
	defer c.Close()
	obs := c.Locals("observer").(*chat.ChatObserver)

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			break
		}

		for _, client := range obs.Clients {
			err := client.WebsocketConn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				break
			}
		}
	}
}

func ShowClientsHandler(c *fiber.Ctx) error {
	obs := c.Locals("observer").(*chat.ChatObserver)
	var clients []chat.ClientJson
	for _, client := range obs.Clients {
		clients = append(clients, chat.ClientJson{Id: client.Id, Name: client.Name})
	}

	return c.JSON(clients)
}

func WsTestHandler(c *websocket.Conn) {
	var (
		msg []byte
		err error
	)
	for {
		if _, msg, err = c.ReadMessage(); err != nil {
			log.Println("read: ", msg)
		}

		msg_str := string(msg)

		if strings.ToLower(msg_str) == "hola" {
			toSend := []byte("Hola, ¿Cómo estás?")

			if err = c.WriteMessage(websocket.TextMessage, toSend); err != nil {
				fmt.Println("write: ", err)
				fmt.Println("write: ", toSend)
			}
		}

		if strings.ToLower(msg_str) == "¿cómo estás?" {
			toSend := []byte("Muy bien, soy un bot tonto")

			if err = c.WriteMessage(websocket.TextMessage, toSend); err != nil {
				fmt.Println("write: ", err)
				fmt.Println("write: ", toSend)
			}
		}

		if strings.ToLower(msg_str) == "chau" {
			toSend := []byte("Hasta pronto, espero haberte ayudado")

			if err = c.WriteMessage(websocket.TextMessage, toSend); err != nil {
				fmt.Println("write: ", err)
				fmt.Println("write: ", toSend)
			}
		}
	}
}
