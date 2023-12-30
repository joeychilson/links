package server

import (
	"net/http"

	"github.com/joeychilson/links/pages/login"
)

func (s *Server) LoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login.Page().Render(r.Context(), w)
	}
}
