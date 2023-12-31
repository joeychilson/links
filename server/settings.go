package server

import (
	"net/http"

	"github.com/joeychilson/links/pages/settings"
)

func (s *Server) SettingsPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := s.UserFromContext(ctx)
		settings.Page(settings.Props{User: user}).Render(ctx, w)
	}
}
