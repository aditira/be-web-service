package module

import (
	"net/http"

	"tira.com/src/helper"
)

// Get userPublic
func UserPublic(w http.ResponseWriter, r *http.Request) {
	helper.CheckMethod(w, r, "GET")

	w.Write([]byte("User Public"))
}

// Get userGuest
func UserGuest(w http.ResponseWriter, r *http.Request) {
	helper.CheckMethod(w, r, "GET")

	w.Write([]byte("User Guest"))
}

// Get userModerator
func UserModerator(w http.ResponseWriter, r *http.Request) {
	helper.CheckMethod(w, r, "GET")

	w.Write([]byte("User Moderator"))
}

// Get userAdmin
func UserAdmin(w http.ResponseWriter, r *http.Request) {
	helper.CheckMethod(w, r, "GET")

	w.Write([]byte("User Admin"))
}
