package routes

import (
	"log"
	"sync"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

type Player struct {
	UserId string
	Name   string
	Status string
}

type ChatMessage struct {
	UserId string
	Name   string
	Msg    string
	Time   string
}

type Room struct {
	RoomId     string
	Name       string
	MapId      string
	MapName    string
	MapImage   string
	Difficulty float64
	MaxPlayer  float64
	MinPlayer  float64
	Owner      string
	Players    map[string]Player
	Chat       []ChatMessage
	mu         sync.Mutex
}

var rooms = make(map[string]*Room)
var roomsMu sync.Mutex

func SocketServerRoute() *socketio.Server {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/online_room", "create_room", func(s socketio.Conn, data map[string]interface{}) {

		roomId := data["roomId"].(string)
		roomName := data["roomName"].(string)

		mapId := data["mapId"].(string)
		mapName := data["mapName"].(string)
		mapImage := data["mapImage"].(string)
		difficulty := data["difficulty"].(float64) // Extract as float64 due to JSON number format

		userId := data["userId"].(string)
		characterName := data["characterName"].(string)

		maxPlayer := data["maxPlayer"].(float64) // Extract as float64 due to JSON number format
		minPlayer := data["minPlayer"].(float64) // Extract as float64 due to JSON number format

		roomsMu.Lock()
		rooms[roomId] = &Room{
			RoomId:     roomId,
			Name:       roomName,
			MapId:      mapId,
			MapName:    mapName,
			MapImage:   mapImage,
			Difficulty: difficulty,
			MaxPlayer:  maxPlayer,
			MinPlayer:  minPlayer,
			Owner:      userId,
			Players: map[string]Player{
				userId: {UserId: userId, Name: characterName, Status: "Not Ready"},
			},
			Chat: []ChatMessage{},
		}
		roomsMu.Unlock()

		log.Println("created room:", roomId, roomName, userId, characterName)
		s.Emit("room_created", map[string]interface{}{
			"roomId":   roomId,
			"roomName": roomName,
			"owner":    userId,
		})
	})

	server.OnEvent("/online_room", "join_room", func(s socketio.Conn, data map[string]interface{}) {
		roomId := data["roomId"].(string)
		userId := data["userId"].(string)
		characterName := data["characterName"].(string)

		roomsMu.Lock()
		room, exists := rooms[roomId]
		if exists {
			room.mu.Lock()
			room.Players[userId] = Player{
				UserId: userId,
				Name:   characterName,
				Status: "Not Ready"}
			room.mu.Unlock()
		}
		roomsMu.Unlock()

		if exists {
			s.Join(roomId)
			log.Printf("%s joined room: %s", userId, roomId)
			server.BroadcastToRoom("/online_room", roomId, "room_detail", map[string]interface{}{
				"id":         room.RoomId,
				"name":       room.Name,
				"mapId":      room.MapId,
				"mapName":    room.MapName,
				"mapImage":   room.MapImage,
				"difficulty": room.Difficulty,
				"maxPlayer":  room.MaxPlayer,
				"minPlayer":  room.MinPlayer,
				"owner":      room.Owner,
			})
			server.BroadcastToRoom("/online_room", roomId, "update_players", map[string]interface{}{
				"players": room.Players,
				"owner":   room.Owner,
			})
			message := ChatMessage{
				UserId: userId,
				Name:   characterName,
				Msg:    "join room",
				Time:   time.Now().Format(time.RFC3339),
			}
			room.mu.Lock()
			room.Chat = append(room.Chat, message)
			room.mu.Unlock()

			server.BroadcastToRoom("/online_room", roomId, "chat_client", message)
		}
	})

	server.OnEvent("/online_room", "update_ready_status", func(s socketio.Conn, data map[string]interface{}) {
		roomId := data["roomId"].(string)
		userId := data["userId"].(string)
		readyStatus := data["readyStatus"].(bool)

		roomsMu.Lock()
		room, exists := rooms[roomId]
		if exists {
			room.mu.Lock()
			if player, ok := room.Players[userId]; ok {
				player.Status = "Ready"
				if !readyStatus {
					player.Status = "Not Ready"
				}
				room.Players[userId] = player
			}
			room.mu.Unlock()
		}
		roomsMu.Unlock()

		if exists {
			server.BroadcastToRoom("/online_room", roomId, "update_players", map[string]interface{}{
				"players": room.Players,
				"owner":   room.Owner,
			})
		}
	})

	server.OnEvent("/online_room", "send_chat_message", func(s socketio.Conn, data map[string]interface{}) {
		roomId := data["roomId"].(string)
		userId := data["userId"].(string)
		characterName := data["characterName"].(string)
		msg := data["msg"].(string)

		message := ChatMessage{
			UserId: userId,
			Name:   characterName,
			Msg:    msg,
			Time:   time.Now().Format(time.RFC3339),
		}

		roomsMu.Lock()
		room, exists := rooms[roomId]
		if exists {
			room.mu.Lock()
			room.Chat = append(room.Chat, message)
			room.mu.Unlock()
		}
		roomsMu.Unlock()

		if exists {
			server.BroadcastToRoom("/online_room", roomId, "chat_client", message)
		}
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
