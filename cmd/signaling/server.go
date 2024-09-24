package signaling

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

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

// HandleSignaling manages WebSocket connections, handles messages, and room lifecycle.
func HandleSignaling(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v\n", err)
		return
	}

	defer func() {
		conn.Close()
		log.Println("WebSocket connection closed")
		handleDisconnect(conn)
	}()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v\n", err)
			return
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(p, &msg); err != nil {
			log.Printf("JSON unmarshal error: %v\n", err)
			continue
		}

		roomID, ok := msg["room"].(string)
		if !ok {
			log.Println("Room ID is missing or invalid")
			conn.WriteJSON(map[string]string{"type": "error", "message": "Invalid or missing room ID"})
			continue
		}

		msgType, ok := msg["type"].(string)
		if !ok {
			log.Println("Message type is missing or invalid")
			conn.WriteJSON(map[string]string{"type": "error", "message": "Invalid or missing message type"})
			continue
		}

		switch msgType {
		case "create":
			handleCreate(conn, roomID)
		case "join":
			handleJoin(conn, roomID)
		case "offer", "answer", "ice-candidate":
			handleSignaling(conn, msg)
		default:
			log.Printf("Unknown message type: %s\n", msgType)
			conn.WriteJSON(map[string]string{"type": "error", "message": "Unknown message type"})
		}
	}
}

// handleCreate handles room creation by the sender.
func handleCreate(conn *websocket.Conn, roomID string) {
	roomsMu.Lock()
	defer roomsMu.Unlock()

	if _, exists := rooms[roomID]; exists {
		conn.WriteJSON(map[string]string{"type": "error", "message": "Room already exists"})
		return
	}

	room := &Room{Sender: conn}
	rooms[roomID] = room
	log.Printf("Room created: %s", roomID)

	response := map[string]string{"type": "room-created", "room": roomID}
	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Error sending room-created response: %v\n", err)
	}
}

// handleJoin allows a receiver to join a room.
func handleJoin(conn *websocket.Conn, roomID string) {
	log.Printf("Attempting to join room: %s", roomID)

	roomsMu.Lock()
	room, exists := rooms[roomID]
	roomsMu.Unlock()

	if !exists {
		log.Printf("Room not found: %s", roomID)
		conn.WriteJSON(map[string]string{"type": "error", "message": "Room not found"})
		return
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	if room.Receiver != nil {
		log.Printf("Room is full: %s", roomID)
		conn.WriteJSON(map[string]string{"type": "error", "message": "Room is full"})
		return
	}

	room.Receiver = conn
	log.Printf("Peer joined room: %s", roomID)

	if err := room.Sender.WriteJSON(map[string]string{"type": "peer-joined"}); err != nil {
		log.Printf("Error notifying sender about peer join: %v\n", err)
	}

	if err := conn.WriteJSON(map[string]string{"type": "joined"}); err != nil {
		log.Printf("Error sending joined response: %v\n", err)
	}
}

// handleSignaling forwards signaling messages between peers in the room.
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
			log.Printf("Error forwarding message: %v\n", err)
		}
	} else {
		log.Println("Target peer not available")
	}
}

// handleDisconnect cleans up the room if one of the peers disconnects.
func handleDisconnect(conn *websocket.Conn) {
	roomsMu.Lock()
	defer roomsMu.Unlock()

	for roomID, room := range rooms {
		room.mu.Lock()

		// Check if the disconnecting connection is the sender or receiver
		if conn == room.Sender {
			log.Printf("Sender left room: %s", roomID)
			room.Sender = nil
		} else if conn == room.Receiver {
			log.Printf("Receiver left room: %s", roomID)
			room.Receiver = nil
		}

		// If both the sender and receiver are nil, remove the room
		if room.Sender == nil && room.Receiver == nil {
			delete(rooms, roomID)
			log.Printf("Room %s deleted", roomID)
		}

		room.mu.Unlock()
	}
}
