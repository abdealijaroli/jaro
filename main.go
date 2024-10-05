package main

import (
	"log"
	"net/http"
	"os"

	"github.com/abdealijaroli/jaro/api"
	"github.com/abdealijaroli/jaro/cmd"
	"github.com/abdealijaroli/jaro/cmd/stream"
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

	http.HandleFunc("/", api.HandleGetAndRedirect(storage))

	http.HandleFunc("/ws", stream.InitiateTransfer(storage))

	if len(os.Args) > 1 {
		cmd.Execute()
		return
	}

	log.Println("Server starting on :8008")
	if err := http.ListenAndServe(":8008", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
