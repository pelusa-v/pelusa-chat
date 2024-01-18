package main

import (
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
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

	app.Post("/register", websocket.New(registerHandler))
	app.Get("/clients", showClientsHandler)
	app.Get("/ws_test", websocket.New(wsTestHandler))
	app.Listen(":3000")
}

func registerHandler(c *websocket.Conn) {
	obs := c.Locals("observer").(*chat.Observer)
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
		mt      int
		msg     []byte
		err     error
		summary string
	)
	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Println("read: ", msg)
		}

		// summary = "Hereeeeee"
		summary = fmt.Sprint(mt) + " / " + string(msg)

		// if err = c.WriteMessage(mt, msg); err != nil {
		toSend := []byte(summary)
		if err = c.WriteMessage(websocket.TextMessage, toSend); err != nil {
			fmt.Println("write: ", err)
			fmt.Println("write: ", toSend)
		}
	}
}
