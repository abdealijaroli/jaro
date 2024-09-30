package main

import (
	"log"
	"net/http"
	"os"

	"github.com/abdealijaroli/jaro/cmd"
	"github.com/abdealijaroli/jaro/cmd/stream"
	"github.com/abdealijaroli/jaro/store"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for this example
	},
}

func main() {
	storage, err := store.NewPostgresStore()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer storage.Close()

	if err := storage.Init(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	fs := http.FileServer(http.Dir("web/public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		shortCode := r.URL.Path[1:]
		if shortCode == "" {
			http.ServeFile(w, r, "web/index.html")
			return
		}
		originalURL, err := storage.GetOriginalURL(shortCode)
		if err != nil {
			http.Error(w, "four oh four - not found :(", http.StatusNotFound)
			return
		}
		isFileTransfer, err := storage.CheckFileTransfer(shortCode)
		if err != nil {
			http.Error(w, "four oh four - not found :(", http.StatusNotFound)
			return
		}
		if isFileTransfer {
			http.ServeFile(w, r, "web/receiver.html")
			return
		}
		http.Redirect(w, r, originalURL, http.StatusFound)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		roomID := r.URL.Query().Get("room")
		if roomID == "" {
			http.Error(w, "Missing room ID", http.StatusBadRequest)
			return
		}

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Error upgrading to WebSocket: %v", err)
			return
		}
		defer conn.Close()

		//Retrieve the file path associated with this room ID from your database
		filePath, err := storage.GetOriginalURL(roomID)
		if err != nil {
			log.Printf("Error getting file path for room %s: %v", roomID, err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error: File not found"))
			return
		}

		// Initiate the file transfer
		stream.HandleTransferRequest(conn, filePath)
	})

	if len(os.Args) > 1 {
		cmd.Execute()
	}

	log.Println("Server starting on :8008")
	if err := http.ListenAndServe(":8008", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
