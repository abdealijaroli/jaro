package api

import (
	"net/http"

	"github.com/abdealijaroli/jaro/store"
)

func AddUserToWaitlist(w http.ResponseWriter, r *http.Request, storage *store.PostgresStore) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	err = storage.CreateWaitlist(name, email)
	if err != nil {
		http.Error(w, "Error adding user to waitlist", http.StatusInternalServerError)
		return
	}
}
