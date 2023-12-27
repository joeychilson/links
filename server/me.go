package server

import (
	"log"
	"net/http"

	"github.com/joeychilson/links/database"
	account "github.com/joeychilson/links/pages/me"
)

func (s *Server) MePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := s.UserFromContext(r.Context())

		likedLinks, err := s.queries.LikedFeed(r.Context(), database.LikedFeedParams{
			UserID: user.ID,
			Limit:  25,
			Offset: 0,
		})
		if err != nil {
			log.Printf("error getting liked links: %v", err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		account.Page(account.Props{User: user, Links: likedLinks}).Render(r.Context(), w)
	}
}
