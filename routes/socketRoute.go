package routes

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

func SocketServerRoute() *socketio.Server {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		log.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		log.Println("chat msg received:", msg)
		s.Emit("reply", "Message received: "+msg)
		return "recv " + msg
	})

	server.OnEvent("/chat", "data", func(s socketio.Conn, data map[string]interface{}) {
		log.Println("Data received:", data)
		response := map[string]interface{}{
			"status":   "success",
			"received": data,
		}
		s.Emit("dataResponse", response)
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		log.Println("closed", msg)
	})

	return server
}
