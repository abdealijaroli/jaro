package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	helloComponent := hello("aj")

	http.Handle("/", templ.Handler(helloComponent))
	fmt.Println("listening on :8008")
	http.ListenAndServe(":8008", nil)
}
