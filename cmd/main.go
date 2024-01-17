package main

import (
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/pelusa-v/pelusa-chat.git/internal/chat"
)

func main() {
	app := fiber.New()

	obs := chat.NewObserver()
	go obs.Start()

	// registerHandlerWithParams := func(obs *chat.Observer) fiber.Handler {
	// 	return func(c *fiber.Ctx) error {
	// 		sampleIds := []string{
	// 			"697b1579-2186-472f-b636-cfe1a2559bc9",
	// 			"a19e7ea4-5b09-45fb-b37b-358ebe0e5aa3",
	// 			"159dce89-dc23-498b-a941-069b7dbbd577",
	// 			"a9f7b966-8652-477c-9139-a14ca5a19669",
	// 		}
	// 		names := []string{
	// 			"Bob",
	// 			"Jorge",
	// 			"Tomy",
	// 			"Li",
	// 		}

	// 		id := sampleIds[rand.Intn(len(sampleIds))]
	// 		name := names[rand.Intn(len(names))]
	// 		newClient := chat.NewClient(id, name, obs)
	// 		newClient.Observer.SubscribeClientChan <- newClient

	// 		return c.SendString("A client was added")
	// 	}
	// }

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("observer", obs)
		return c.Next()
	})

	app.Post("/register", registerHandler)
	app.Get("/clients", showClientsHandler)
	app.Listen(":3000")
}

func registerHandler(c *fiber.Ctx) error {
	sampleIds := []string{
		"697b1579-2186-472f-b636-cfe1a2559bc9",
		"a19e7ea4-5b09-45fb-b37b-358ebe0e5aa3",
		"159dce89-dc23-498b-a941-069b7dbbd577",
		"a9f7b966-8652-477c-9139-a14ca5a19669",
	}
	names := []string{
		"Bob",
		"Jorge",
		"Tomy",
		"Li",
	}

	obs := c.Locals("observer").(*chat.Observer)
	id := sampleIds[rand.Intn(len(sampleIds))]
	name := names[rand.Intn(len(names))]
	newClient := chat.NewClient(id, name, obs)
	newClient.Observer.SubscribeClientChan <- newClient

	return c.SendString("A client was added")
}

func showClientsHandler(c *fiber.Ctx) error {
	obs := c.Locals("observer").(*chat.Observer)
	var clients []chat.ClientJson
	for _, client := range obs.Clients {
		clients = append(clients, chat.ClientJson{Id: client.Id, Name: client.Name})
	}

	return c.JSON(clients)
}
