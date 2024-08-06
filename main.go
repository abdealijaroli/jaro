package main

import (
	"fmt"
	"net/http"

	// "github.com/a-h/templ"
)

func main() {
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	
	

	fmt.Println("listening on :8008")
	err := http.ListenAndServe(":8008", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
