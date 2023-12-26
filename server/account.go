package server

import (
	"net/http"

	"github.com/joeychilson/links/templates/pages/account"
)

func (s *Server) AccountPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := s.UserFromContext(r.Context())
		account.Page(account.Props{User: user}).Render(r.Context(), w)
	}
}
