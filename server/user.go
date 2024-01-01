package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"

	"github.com/joeychilson/links/pages/user"
)

func (s *Server) UserPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		u := s.UserFromContext(ctx)
		username := chi.URLParam(r, "username")

		userProfile, err := s.queries.UserProfile(ctx, username)
		if err != nil {
			oplog.Error("error getting user profile", err)
			s.Redirect(w, r, "/")
			return
		}

		user.Page(user.Props{User: u, Profile: userProfile}).Render(ctx, w)
	}
}
