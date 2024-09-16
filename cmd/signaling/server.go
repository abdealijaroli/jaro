package signaling

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Room struct {
	Sender   *websocket.Conn
	Receiver *websocket.Conn
	mu       sync.Mutex
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	rooms   = make(map[string]*Room)
	roomsMu sync.Mutex
)

func StartSignalingServer() {
	http.HandleFunc("/ws", HandleWebSocket)
	log.Fatal(http.ListenAndServe(":8008", nil))
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(p, &msg); err != nil {
			log.Println(err)
			continue
		}

		_, ok := msg["room"].(string)
		if !ok {
			log.Println("Room ID is missing or not a string")
			continue
		}

		switch msg["type"] {
		case "create":
			handleCreate(conn)
		case "join":
			handleJoin(conn, msg["room"].(string))
		case "offer", "answer", "ice-candidate":
			handleSignaling(conn, msg)
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func handleCreate(conn *websocket.Conn) {
	roomID := uuid.New().String()
	room := &Room{Sender: conn}

	roomsMu.Lock()
	rooms[roomID] = room
	roomsMu.Unlock()

	response := map[string]string{"type": "room-created", "room": roomID}
	if err := conn.WriteJSON(response); err != nil {
		log.Println(err)
	}
}

func handleJoin(conn *websocket.Conn, roomID string) {
	roomsMu.Lock()
	room, exists := rooms[roomID]
	roomsMu.Unlock()

	if !exists {
		conn.WriteJSON(map[string]string{"type": "error", "message": "Room not found"})
		return
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	if room.Receiver != nil {
		conn.WriteJSON(map[string]string{"type": "error", "message": "Room is full"})
		return
	}

	room.Receiver = conn
	conn.WriteJSON(map[string]string{"type": "joined"})
	room.Sender.WriteJSON(map[string]string{"type": "peer-joined"})
}

func handleSignaling(conn *websocket.Conn, msg map[string]interface{}) {
	roomID := msg["room"].(string)
	roomsMu.Lock()
	room, exists := rooms[roomID]
	roomsMu.Unlock()

	if !exists {
		conn.WriteJSON(map[string]string{"type": "error", "message": "Room not found"})
		return
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	var target *websocket.Conn
	if conn == room.Sender {
		target = room.Receiver
	} else {
		target = room.Sender
	}

	if target != nil {
		if err := target.WriteJSON(msg); err != nil {
			log.Println(err)
		}
	}
}
