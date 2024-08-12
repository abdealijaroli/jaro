package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"

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

	cmd.Execute()
}
