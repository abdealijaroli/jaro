package stream

import (
	"encoding/binary"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/pion/webrtc/v3"

	"github.com/abdealijaroli/jaro/cmd/webrtcconn"
)

func ReceiveFile(roomID string) {
	peerConnection := webrtcconn.CreatePeerConnection()

	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		log.Println("New DataChannel:", d.Label())
		receiveFile(d)
	})

	webrtcconn.HandleSignaling(peerConnection, roomID)
}

var highestChunkReceived uint32

func receiveFile(dataChannel *webrtc.DataChannel) {
	var file *os.File
	var fileSize int64
	var receivedSize int64
	var chunks map[uint32][]byte
	var mu sync.Mutex

	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		// Handle Metadata (initial file info)
		if file == nil {
			metadata := string(msg.Data)
			parts := strings.Split(metadata, ":")
			if len(parts) != 2 {
				log.Println("Invalid metadata")
				return
			}

			name := parts[0]
			size, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				log.Println("Failed to parse file size:", err)
				return
			}
			file, err = os.Create(name)
			if err != nil {
				log.Println("Failed to create file:", err)
				return
			}

			fileSize = size
			chunks = make(map[uint32][]byte)
			log.Printf("Receiving file: %s, Size: %d bytes\n", name, fileSize)
			return
		}

		// Handle EOF
		if string(msg.Data) == "EOF" {
			log.Println("End of File signal received, reassembling...")
			reassembleFile(file, chunks)
			return
		}

		// Handle file chunks
		chunkNumber := binary.BigEndian.Uint32(msg.Data[:4])
		data := msg.Data[4:]

		mu.Lock()
		chunks[chunkNumber] = data
		receivedSize += int64(len(data))
		if chunkNumber > highestChunkReceived {
			highestChunkReceived = chunkNumber
		}
		mu.Unlock()

		// Log progress
		log.Printf("Chunk %d received, total received size: %d/%d bytes\n", chunkNumber, receivedSize, fileSize)

		// Close the file if we've received all data
		if receivedSize >= fileSize {
			file.Close()
			log.Println("File transfer complete, file closed.")
		}
	})
}

func reassembleFile(file *os.File, chunks map[uint32][]byte) {
	for i := uint32(0); i <= highestChunkReceived; i++ {
		if chunk, exists := chunks[i]; exists {
			_, err := file.Write(chunk)
			if err != nil {
				log.Println("Error writing chunk:", err)
				return
			}
		} else {
			log.Printf("Missing chunk %d, skipping...\n", i)
		}
	}
	file.Close()
	log.Println("File reassembled successfully.")
}
