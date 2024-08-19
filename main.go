package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"

	"github.com/abdealijaroli/jaro/api"
	"github.com/abdealijaroli/jaro/cmd"
	"github.com/abdealijaroli/jaro/store"
	"github.com/abdealijaroli/jaro/web/components"
)

func main() {
	// db init
	storage, err := store.NewPostgresStore()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer storage.Close()

	err = storage.Init()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	http.Handle("/hello", templ.Handler(components.Hello("abdeali", "click me")))

	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		api.AddUserToWaitlist(w, r, storage)
	})
	
	cmd.Execute()

	fmt.Println("Server is running on port 8008")
	http.ListenAndServe(":8008", nil)
}
