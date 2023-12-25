package server

import (
	"net/http"

	"github.com/joeychilson/lixy/models"
	"github.com/joeychilson/lixy/pages/account"
)

func (s *Server) handleAccountPage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userKey).(*models.User)
	if user == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	account.Page(account.Props{User: user}).Render(r.Context(), w)
}
