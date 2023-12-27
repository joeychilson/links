package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/joeychilson/links/database"
	"github.com/joeychilson/links/pages/link"
)

func (s *Server) Link() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := s.UserFromContext(r.Context())
		linkID := r.URL.Query().Get("id")

		if linkID == "" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		linkUUID, err := uuid.Parse(linkID)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		var userID uuid.UUID
		if user != nil {
			userID = user.ID
		} else {
			userID = uuid.Nil
		}

		linkRow, err := s.queries.Link(r.Context(), database.LinkParams{
			UserID: userID,
			LinkID: linkUUID,
		})
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		link.Page(link.Props{User: user, Link: linkRow}).Render(r.Context(), w)
	}
}
