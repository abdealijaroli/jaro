package main

import (
	"fmt"
	"net/http"

	// "github.com/a-h/templ"
)

func main() {
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	// http.Handle("/", templ.Handler(hello(`hello`, `aj`)))

	fmt.Println("listening on :8008")
	http.ListenAndServe(":8008", nil) 
}