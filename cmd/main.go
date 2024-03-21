package main

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/pelusa-v/pelusa-chat.git/internal/chat"
	"github.com/pelusa-v/pelusa-chat.git/internal/handlers"
)

func main() {
	app := fiber.New()

	go chat.Manager.Start()

	app.Get("/api/ws/register/:nick", websocket.New(handlers.RegisterHandler))
	// app.Get("/api/clients", handlers.ShowClientsHandler)

	app.Get("/room/:nick", handlers.ChatRoomViewHandler)
	app.All("/", handlers.RegisterRoomViewHandler)

	app.Listen("127.0.0.1:3000")
}
