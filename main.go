package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	var helloComponent templ.Component
	yes := true
	if yes {
		helloComponent = hello("aj", "click me yes")
	} else {
		helloComponent = hello("aj", "no click me")
	}

	http.Handle("/", templ.Handler(helloComponent))

	fmt.Println("listening on :8008")
	http.ListenAndServe(":8008", nil)
}
