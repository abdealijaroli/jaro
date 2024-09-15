package stream

import (
	"encoding/binary"
	"fmt"
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
		fmt.Println("New DataChannel:", d.Label())
		receiveFile(d)
	})

	webrtcconn.HandleSignaling(peerConnection, roomID)
}

var highestChunkReceived uint32

func receiveFile(dataChannel *webrtc.DataChannel) {
	var file *os.File
	var fileSize int64 // needed? how?
	var receivedSize int64
	var chunks map[uint32][]byte
	var mu sync.Mutex

	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		if file == nil {
			metadata := string(msg.Data)
			parts := strings.Split(metadata, ":")
			if len(parts) != 2 {
				fmt.Println("Invalid metadata")
				return
			}

			name := parts[0]
			size, _ := strconv.ParseInt(parts[1], 10, 64)
			file, _ = os.Create(name)
			fileSize = size
			chunks = make(map[uint32][]byte)
		}

		if string(msg.Data) == "EOF" {
			reassembleFile(file, chunks)
			return
		}

		chunkNumber := binary.BigEndian.Uint32(msg.Data[:4])
		data := msg.Data[4:]
		mu.Lock()
		chunks[chunkNumber] = data
		if chunkNumber > highestChunkReceived {
			highestChunkReceived = chunkNumber
		}
		mu.Unlock()

		if receivedSize >= fileSize {
			file.Close()
		}
	})

}

func reassembleFile(file *os.File, chunks map[uint32][]byte) {
	for i := uint32(0); i <= highestChunkReceived; i++ {
		if chunk, exists := chunks[i]; exists {
			file.Write(chunk)
		}
	}

}
