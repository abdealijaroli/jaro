package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/abdealijaroli/jaro/web/components"

	"github.com/abdealijaroli/jaro/cmd"
)

func main() {
	// db init
	store, err := NewPostgresStore()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer store.db.Close()

	err = store.Init()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
 
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	http.Handle("/hello", templ.Handler(components.Hello("abdeali", "click me")))

	// cli init
	cmd.Execute()

	fmt.Println("listening on :8008")
	err = http.ListenAndServe(":8008", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
