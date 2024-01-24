package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pelusa-v/pelusa-chat.git/internal/chat"
)

func main() {
	app := fiber.New()

	obs := chat.NewObserver()
	go obs.Start()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("observer", obs)
		return c.Next()
	})

	app.Get("/ws/register/:nick", websocket.New(registerHandler))
	app.Get("/clients", showClientsHandler)
	app.Get("/room/:nick", chatRoomHandler)
	app.All("/register", registerRoomHandler)
	// app.Get("/ws_test", websocket.New(wsTestHandler))
	// app.Get("/we_broadcast", websocket.New(broacastHandler))
	app.Listen(":3000")
}

func registerRoomHandler(c *fiber.Ctx) error {

	if c.Method() == fiber.MethodPost {
		nickName := c.FormValue("nick")
		return c.Redirect(fmt.Sprintf("/room/%s", nickName))
	}

	return c.Render("internal/views/register.html", nil)
}

func chatRoomHandler(c *fiber.Ctx) error {
	data := fiber.Map{
		"nick": c.Params("nick"),
	}
	return c.Render("internal/views/room.html", data)
}

func registerHandler(c *websocket.Conn) {
	obs := c.Locals("observer").(*chat.Observer)

	client := chat.NewClient(uuid.New().String(), c.Params("nick"), obs, c)
	client.Observer.SubscribeClientChan <- client

	for {
		_, msg, _ := client.WebsocketConn.ReadMessage()
		fmt.Println(msg)
		// _, msg, err := c.ReadMessage()
		// genericResponse := []byte("Respuesta genérica a :" + string(msg))
		// client.WebsocketConn.WriteMessage(websocket.TextMessage, genericResponse)
	}
}

func broacastHandler(c *websocket.Conn) {
	defer c.Close()
	obs := c.Locals("observer").(*chat.Observer)

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

func readWebSocketMessage() {

}

// func registerHandler(c *fiber.Ctx) error {
// 	sampleIds := []string{
// 		"697b1579-2186-472f-b636-cfe1a2559bc9",
// 		"a19e7ea4-5b09-45fb-b37b-358ebe0e5aa3",
// 		"159dce89-dc23-498b-a941-069b7dbbd577",
// 		"a9f7b966-8652-477c-9139-a14ca5a19669",
// 	}
// 	names := []string{
// 		"Bob",
// 		"Jorge",
// 		"Tomy",
// 		"Li",
// 	}

// 	obs := c.Locals("observer").(*chat.Observer)
// 	id := sampleIds[rand.Intn(len(sampleIds))]
// 	name := names[rand.Intn(len(names))]
// 	newClient := chat.NewClient(id, name, obs)
// 	newClient.Observer.SubscribeClientChan <- newClient

// 	return c.SendString("A client was added")
// }

func showClientsHandler(c *fiber.Ctx) error {
	obs := c.Locals("observer").(*chat.Observer)
	var clients []chat.ClientJson
	for _, client := range obs.Clients {
		clients = append(clients, chat.ClientJson{Id: client.Id, Name: client.Name})
	}

	return c.JSON(clients)
}

func wsTestHandler(c *websocket.Conn) {
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
