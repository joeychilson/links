package server

import (
	"net/http"

	"github.com/go-chi/httplog/v2"
)

func (s *Server) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		if user == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		err := s.sessionManager.Delete(w, r)
		if err != nil {
			oplog.Error("failed to delete session", "error", err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		oplog.Info("user logged out", "user_id", user.ID.String())
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
