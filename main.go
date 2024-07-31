package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	http.Handle("/", templ.Handler(hello(`hello`, `aj`)))

	fmt.Println("listening on :8008")
	http.ListenAndServe(":8008", nil)
}
