package server

import (
	"net/http"
)

func (s *Server) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := s.UserFromContext(r.Context())
		if user == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Delete session
		err := s.sessionManager.Delete(w, r)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Redirect to the home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
