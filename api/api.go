package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/abdealijaroli/jaro/store"
)



func AddUserToWaitlist(w http.ResponseWriter, r *http.Request, storage *store.PostgresStore) {
	var name, email string

	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
		var data struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			log.Printf("JSON decode error: %v\n", err)
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}
		name, email = data.Name, data.Email
	} else {
		if err := r.ParseForm(); err != nil {
			log.Printf("Form parse error: %v\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		name = r.FormValue("name")
		email = r.FormValue("email")
	}

	if name == "" || email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	err := storage.CreateWaitlist(name, email)
	if err != nil {
		log.Printf("CreateWaitlist error: %v\n", err)
		http.Error(w, "Error adding to waitlist", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("You are on the waitlist!"))
}
