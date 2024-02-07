package main

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/pelusa-v/pelusa-chat.git/internal/chat"
	"github.com/pelusa-v/pelusa-chat.git/internal/handlers"
)

func main() {
	app := fiber.New()

	obs := chat.NewChatObserver()
	go obs.Start()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("observer", obs)
		return c.Next()
	})

	app.Get("/api/ws/register/:nick", websocket.New(handlers.RegisterHandler))
	app.Get("/api/clients", handlers.ShowClientsHandler)
	app.Get("/room/:nick", handlers.ChatRoomHandler)
	app.All("/", handlers.RegisterRoomHandler)
	// app.Get("/ws_test", websocket.New(wsTestHandler))
	// app.Get("/we_broadcast", websocket.New(broacastHandler))
	app.Listen(":3000")
}
