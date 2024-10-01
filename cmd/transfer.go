package cmd

import (
	"log"

	"github.com/abdealijaroli/jaro/cmd/utils"
	"github.com/abdealijaroli/jaro/cmd/webrtcconn"
	"github.com/abdealijaroli/jaro/store"
)

func TransferFile(filePath string, store *store.PostgresStore) {
    shortURL, roomID := utils.GenerateShortCode(filePath)
    utils.GenerateQRCode(shortURL)

    pc := webrtcconn.CreatePeerConnection()
    // log.Printf("Peer connection created for %s\n", roomID)

    offer, err := pc.CreateOffer(nil)
    if err != nil {
        log.Fatalf("Failed to create offer: %v", err)
    }

    err = pc.SetLocalDescription(offer)
    if err != nil {
        log.Fatalf("Failed to set local description: %v", err)
    }

    // log.Printf("Local description set for %s\n", roomID)

    // Handle signaling messages
    go webrtcconn.HandleSignaling(pc, roomID)
    // stream.InitiateTransfer(conn, filePath, roomID, store)
}