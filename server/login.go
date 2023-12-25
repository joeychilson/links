package server

import (
	"net/http"

	"github.com/joeychilson/inquire/pages/login"
)

func (s *Server) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	login.Page().Render(r.Context(), w)
}
