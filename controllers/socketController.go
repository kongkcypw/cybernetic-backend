package controllers

import (
	"fmt"

	"github.com/gofiber/contrib/socketio"
)

func SocketController(kws *socketio.Websocket) {

	// The key for the map is message.to
	clients := make(map[string]string)

	// Retrieve the user id from endpoint
	userId := kws.Params("id")

	// Add the connection to the list of the connected clients
	// The UUID is generated randomly and is the key that allow
	// socketio to manage Emit/EmitTo/Broadcast
	clients[userId] = kws.UUID

	// Every websocket connection has an optional session key => value storage
	kws.SetAttribute("user_id", userId)

	//Broadcast to all the connected users the newcomer
	kws.Broadcast([]byte(fmt.Sprintf("New user connected: %s and UUID: %s", userId, kws.UUID)), true, socketio.TextMessage)
	//Write welcome message
	kws.Emit([]byte(fmt.Sprintf("Hello user: %s with UUID: %s", userId, kws.UUID)), socketio.TextMessage)
}
