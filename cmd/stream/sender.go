package stream

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/abdealijaroli/jaro/store"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// InitiateTransfer starts the file transfer process
func InitiateTransfer(w http.ResponseWriter, r *http.Request, filePath string, roomID string, store *store.PostgresStore) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Printf("Initiating transfer for room: %s\n", roomID)

	// Store file metadata in the database
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Error getting file info: %v\n", err)
		return
	}

	err = store.AddShortURLToDB(fileInfo.Name(), roomID, true)
	if err != nil {
		log.Printf("Error adding short URL to DB: %v\n", err)
		return
	}

	// Start the file transfer
	HandleTransferRequest(conn, filePath)
}

func HandleTransferRequest(conn *websocket.Conn, filePath string) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("File open error: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error opening file"))
		return
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("File stat error: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error getting file info"))
		return
	}

	// Send file metadata
	metadata := map[string]interface{}{
		"name": fileInfo.Name(),
		"size": fileInfo.Size(),
	}
	if err := conn.WriteJSON(metadata); err != nil {
		log.Printf("Error sending file metadata: %v", err)
		return
	}

	// Send file content
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("File read error: %v", err)
			return
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, buffer[:n]); err != nil {
			log.Printf("Error sending file chunk: %v", err)
			return
		}
	}

	log.Println("File sent successfully")
}

func CreateRoom(conn *websocket.Conn, roomID string) error {
	// This function remains largely unchanged
	log.Printf("Creating room: %s", roomID)
	response := map[string]string{"type": "room-created", "room": roomID}
	err := conn.WriteJSON(response)
	if err != nil {
		log.Printf("Error sending room-created response: %v\n", err)
	}
	return err
}

// package stream

// import (
// 	"bufio"
// 	"encoding/binary"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"
// 	"sync"

// 	"github.com/gorilla/websocket"
// 	"github.com/pion/webrtc/v3"

// 	"github.com/abdealijaroli/jaro/cmd/signaling"
// 	// "github.com/abdealijaroli/jaro/cmd/webrtcconn"
// 	"github.com/abdealijaroli/jaro/store"
// )

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// // InitiateTransfer starts the file transfer process and returns the URL for the receiver.
// func InitiateTransfer(conn *websocket.Conn, filePath string, roomID string, store *store.PostgresStore) {
// 	err := CreateRoom(conn, roomID)
// 	log.Printf("Creating room from sender: %s\n", roomID)
// 	if err != nil {
// 		log.Printf("Error creating room: %v\n", err)
// 	}
// 	startFileTransfer(filePath, roomID, store)
// }

// func HandleTransferRequest(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("Upgrade error:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	// Open the file
// 	file, err := os.Open("test_img.jpg")
// 	if err != nil {
// 		log.Println("File open error:", err)
// 		return
// 	}
// 	defer file.Close()

// 	// Get file info
// 	fileInfo, err := file.Stat()
// 	if err != nil {
// 		log.Println("File stat error:", err)
// 		return
// 	}

// 	// Send file name and size
// 	err = conn.WriteJSON(map[string]interface{}{
// 		"name": fileInfo.Name(),
// 		"size": fileInfo.Size(),
// 	})
// 	if err != nil {
// 		log.Println("Send file info error:", err)
// 		return
// 	}

// 	// Send file content
// 	writer, err := conn.NextWriter(websocket.BinaryMessage)
// 	if err != nil {
// 		log.Println("Get writer error:", err)
// 		return
// 	}
// 	if _, err := io.Copy(writer, file); err != nil {
// 		log.Println("File send error:", err)
// 		return
// 	}
// 	if err := writer.Close(); err != nil {
// 		log.Println("Writer close error:", err)
// 		return
// 	}

// 	log.Println("File sent successfully")
// }

// func CreateRoom(conn *websocket.Conn, roomID string) error {
// 	signaling.RoomsMu.Lock()
// 	defer signaling.RoomsMu.Unlock()

// 	_, exists := signaling.Rooms[roomID]
// 	if exists {
// 		log.Printf("Room already exists: %s", roomID)
// 		return fmt.Errorf("room already exists")
// 	}

// 	room := &signaling.Room{Sender: conn}
// 	signaling.Rooms[roomID] = room
// 	log.Printf("Room created from sender: %s", roomID)

// 	response := map[string]string{"type": "room-created", "room": roomID}
// 	err := conn.WriteJSON(response)
// 	if err != nil {
// 		log.Printf("Error sending room-created response: %v\n", err)
// 	}
// 	return err
// }

// // startFileTransfer initiates the file transfer using a WebRTC DataChannel.
// func startFileTransfer(filePath string, roomID string, store *store.PostgresStore) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("Upgrade error:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	file, err := os.Open("test_img.jpg")
// 	if err != nil {
// 		log.Println("File open error:", err)
// 		return
// 	}
// 	defer file.Close()

// 	// Get file info
// 	fileInfo, err := file.Stat()
// 	if err != nil {
// 		log.Println("File stat error:", err)
// 		return
// 	}

// 	// Send file name and size
// 	err = conn.WriteJSON(map[string]interface{}{
// 		"name": fileInfo.Name(),
// 		"size": fileInfo.Size(),
// 	})
// 	if err != nil {
// 		log.Println("Send file info error:", err)
// 		return
// 	}

// 	// Send file content
// 	writer, err := conn.NextWriter(websocket.BinaryMessage)
// 	if err != nil {
// 		log.Println("Get writer error:", err)
// 		return
// 	}
// 	if _, err := io.Copy(writer, file); err != nil {
// 		log.Println("File send error:", err)
// 		return
// 	}
// 	if err := writer.Close(); err != nil {
// 		log.Println("Writer close error:", err)
// 		return
// 	}

// 	log.Println("File sent successfully")
// 	// file, err := os.Open(filePath)
// 	// if err != nil {
// 	// 	fmt.Printf("Error opening file: %v\n", err)
// 	// 	return
// 	// }
// 	// defer file.Close()

// 	// fileInfo, err := file.Stat()
// 	// if err != nil {
// 	// 	fmt.Printf("Error getting file info: %v\n", err)
// 	// 	return
// 	// }

// 	// // Store file metadata (filename and roomID) in the database
// 	// err = store.AddShortURLToDB(fileInfo.Name(), roomID, true)
// 	// if err != nil {
// 	// 	fmt.Printf("Error adding short URL to DB: %v\n", err)
// 	// 	return
// 	// }

// 	// // Create a new peer connection
// 	// peerConnection := webrtcconn.CreatePeerConnection()

// 	// // Create a data channel for the file transfer
// 	// dataChannel, err := peerConnection.CreateDataChannel("fileTransfer", nil)
// 	// if err != nil {
// 	// 	fmt.Printf("Error creating data channel: %v\n", err)
// 	// 	return
// 	// }

// 	// // When the data channel is open, start sending the file
// 	// dataChannel.OnOpen(func() {
// 	// 	log.Println("DataChannel opened")
// 	// 	go sendFile(file, fileInfo, dataChannel)
// 	// })

// 	// // Handle the WebRTC signaling process
// 	// webrtcconn.HandleSignaling(peerConnection, roomID)
// }

// // sendFile sends the file in chunks over the WebRTC DataChannel.
// func sendFile(file *os.File, fileInfo os.FileInfo, dataChannel *webrtc.DataChannel) {
// 	log.Println("Sending file:", fileInfo.Name())
// 	metadata := fmt.Sprintf("%s:%d", fileInfo.Name(), fileInfo.Size())
// 	dataChannel.Send([]byte(metadata))

// 	reader := bufio.NewReader(file)

// 	var wg sync.WaitGroup
// 	chunks := make(chan []byte, 100)

// 	// Start 5 workers to send file chunks concurrently
// 	for i := 0; i < 5; i++ {
// 		go func() {
// 			for chunk := range chunks {
// 				dataChannel.Send(chunk)
// 				wg.Done()
// 			}
// 		}()
// 	}

// 	// Read and send file in chunks
// 	chunkNumber := uint32(0)
// 	for {
// 		chunk := make([]byte, 1024)
// 		n, err := reader.Read(chunk)
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			fmt.Printf("Error reading file: %v\n", err)
// 			return
// 		}

// 		// Add chunk number to the chunk data
// 		chunkWithNumber := make([]byte, 4+n)
// 		binary.BigEndian.PutUint32(chunkWithNumber[:4], chunkNumber)
// 		copy(chunkWithNumber[4:], chunk[:n])

// 		wg.Add(1)
// 		chunks <- chunkWithNumber

// 		chunkNumber++
// 	}

// 	// Close the chunks channel and wait for all chunks to be sent
// 	close(chunks)
// 	wg.Wait()

// 	// Send an EOF signal to indicate that the transfer is complete
// 	dataChannel.Send([]byte("EOF"))
// }
