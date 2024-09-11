package stream

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/pion/webrtc/v3"

	"github.com/abdealijaroli/jaro/cmd/utils"
	"github.com/abdealijaroli/jaro/cmd/webrtcconn"
	"github.com/abdealijaroli/jaro/store"
)

const chunkSize = 16384

func InitiateTransfer(filePath string, store *store.PostgresStore) string {
	roomID := utils.GenerateShortCode(filePath)
	go startFileTransfer(filePath, roomID, store)
	return fmt.Sprintf("https://jaroli.me/%s", roomID)
}

func startFileTransfer(filePath string, roomID string, store *store.PostgresStore) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
		return
	}

	// add file to store
	err = store.AddShortURLToDB(fileInfo.Name(), roomID, true)
	if err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
		return
	}

	peerConnection := webrtcconn.CreatePeerConnection()
	dataChannel, err := peerConnection.CreateDataChannel("fileTransfer", nil)
	if err != nil {
		fmt.Printf("Error creating data channel: %v\n", err)
		return
	}

	dataChannel.OnOpen(func() {
		sendFile(file, fileInfo, dataChannel)
	})

	webrtcconn.HandleSignaling(peerConnection, roomID)
}

func sendFile(file *os.File, fileInfo os.FileInfo, dataChannel *webrtc.DataChannel) {
	metadata := fmt.Sprintf("%s:%d", fileInfo.Name(), fileInfo.Size())
	dataChannel.Send([]byte(metadata))

	reader := bufio.NewReader(file)
	var wg sync.WaitGroup
	chunks := make(chan []byte, 100)

	for i := 0; i < 5; i++ {
		go func() {
			for chunk := range chunks {
				dataChannel.Send(chunk)
				wg.Done()
			}
		}()
	}

	chunkNumber := uint32(0)
	for {
		chunk := make([]byte, chunkSize)
		n, err := reader.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}

		chunkWithNumber := make([]byte, 4+n)
		binary.BigEndian.PutUint32(chunkWithNumber[:4], chunkNumber)
		copy(chunkWithNumber[4:], chunk[:n])

		wg.Add(1)
		chunks <- chunkWithNumber

		chunkNumber++
	}

	close(chunks)
	wg.Wait()

	dataChannel.Send([]byte("EOF"))
}
