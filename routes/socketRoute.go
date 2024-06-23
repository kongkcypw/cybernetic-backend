package routes

import (
	controllers "example/backend/controllers"
	"fmt"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SocketRoute(router *fiber.App) {

	// Enable CORS
	router.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Set("Access-Control-Allow-Headers", "Content-Type")
		c.Set("Access-Control-Allow-Credentials", "true")
		return c.Next()
	})

	// Setup the middleware to retrieve the data sent in first GET request
	router.Use(func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	socketio.On(socketio.EventConnect, func(ep *socketio.EventPayload) {
		fmt.Printf("Connected with ID: %s\n", ep.Kws.UUID)
	})

	// On disconnect event
	socketio.On(socketio.EventDisconnect, func(ep *socketio.EventPayload) {
		fmt.Printf("Disconnect\n")
	})

	// On close event
	// This event is called when the server disconnects the user actively with .Close() method
	socketio.On(socketio.EventClose, func(ep *socketio.EventPayload) {
		fmt.Printf("Close\n")
	})

	// On error event
	socketio.On(socketio.EventError, func(ep *socketio.EventPayload) {
		fmt.Printf("Error: %s\n", ep.Error)
	})

	// Custom event handling supported
	socketio.On("send_message", func(ep *socketio.EventPayload) {
		fmt.Printf("message comming %s\n", string(ep.Data))
	})

	router.Get("/socket.io/test", socketio.New(controllers.SocketController))
}
