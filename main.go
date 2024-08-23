package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"

	// "github.com/abdealijaroli/jaro/auth"
	"github.com/abdealijaroli/jaro/api"
	"github.com/abdealijaroli/jaro/cmd"
	"github.com/abdealijaroli/jaro/store"
	"github.com/abdealijaroli/jaro/web/components"
)

func main() {
	// db init
	storage, err := store.NewPostgresStore()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer storage.Close()

	if err := storage.Init(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	http.Handle("/hello", templ.Handler(components.Hello("hi", "click me")))

	http.HandleFunc("POST /signup", func(w http.ResponseWriter, r *http.Request) {
		api.AddUserToWaitlist(w, r, storage)
	})
	// http.HandleFunc("/shorten", auth.AuthMiddleware(api.ShortenURL))
	// http.HandleFunc("/transfer", auth.AuthMiddleware(api.TransferFile))

	cmd.Execute()

	fmt.Println("Server is running on port 8008")
	http.ListenAndServe(":8008", nil)
}
