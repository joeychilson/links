package server

import (
	"net/http"

	"github.com/joeychilson/lixy/database"
	"github.com/joeychilson/lixy/models"
)

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(userKey).(*models.User)
	if user == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Attempt to retrieve the session cookie
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Delete the session token from the database
	err = s.queries.DeleteUserToken(r.Context(), database.DeleteUserTokenParams{
		Token:   cookie.Value,
		Context: "session",
	})
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Delete the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	})

	// Redirect to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
