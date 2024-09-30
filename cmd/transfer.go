package cmd

import (
	"log"

	"github.com/gorilla/websocket"

	"github.com/abdealijaroli/jaro/cmd/stream"
	"github.com/abdealijaroli/jaro/cmd/utils"
	"github.com/abdealijaroli/jaro/cmd/webrtcconn"
	"github.com/abdealijaroli/jaro/store"
)

func TransferFile(filePath string, store *store.PostgresStore) {
    shortURL, roomID := utils.GenerateShortCode(filePath)
    utils.GenerateQRCode(shortURL)

    pc := webrtcconn.CreatePeerConnection()
    log.Printf("Peer connection created for %s\n", roomID)

    offer, err := pc.CreateOffer(nil)
    if err != nil {
        log.Fatalf("Failed to create offer: %v", err)
    }

    err = pc.SetLocalDescription(offer)
    if err != nil {
        log.Fatalf("Failed to set local description: %v", err)
    }

    log.Printf("Local description set for %s\n", roomID)

    // Establish WebSocket connection for signaling
    conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
    if err != nil {
        log.Fatalf("Failed to connect to signaling server: %v", err)
    }
    defer conn.Close()

    // Create the room on the signaling server
    err = stream.CreateRoom(conn, roomID)
    if err != nil {
        log.Fatalf("Failed to create room: %v", err)
    }

    // Send the offer to the signaling server
    conn.WriteJSON(map[string]string{"type": "offer", "room": roomID, "sdp": offer.SDP})

    log.Printf("Signaling connection established for %s\n", roomID)

    // Handle signaling messages
    go webrtcconn.HandleSignaling(pc, roomID)

	

    // stream.InitiateTransfer(conn, filePath, roomID, store)
}
// func ReceiveFile(roomID string) {
// 	stream.ReceiveFile(roomID)
// }
