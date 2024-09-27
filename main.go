package main

import (
	"log"
	"net/http"

	"os"

	"github.com/abdealijaroli/jaro/cmd"
	"github.com/abdealijaroli/jaro/cmd/signaling"
	"github.com/abdealijaroli/jaro/store"
)

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

	if len(os.Args) > 1 {
		cmd.Execute()
	}

	go func() {
		log.Println("HTTP server running on :8008")
		if err := http.ListenAndServe(":8008", nil); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Start the WebSocket server on port 8080
	http.HandleFunc("/ws", signaling.HandleSignaling)
	log.Println("WebSocket server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("WebSocket server error: %v", err)
	}
}
