package main

import (
	"log"
	"net/http"

	"github.com/abdealijaroli/jaro/cmd"
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

	fs := http.FileServer(http.Dir("web"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

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
		http.Redirect(w, r, originalURL, http.StatusFound)
	})

	cmd.Execute()

	// http.ListenAndServe(":8008", nil)
}
