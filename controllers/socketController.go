package controllers

import (
	"fmt"

	"github.com/gofiber/contrib/socketio"
)

func SocketController(kws *socketio.Websocket) {
	kws.Emit([]byte(fmt.Sprintf("Test emit: %s", kws.UUID)), socketio.TextMessage)
}
