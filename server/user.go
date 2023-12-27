package server

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/joeychilson/links/database"
	userpage "github.com/joeychilson/links/pages/user"
)

func (s *Server) UserPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := s.UserFromContext(r.Context())

		username := r.URL.Query().Get("name")
		if username == "" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		pageStr := r.URL.Query().Get("page")
		if pageStr == "" {
			pageStr = "1"
		}

		page, err := strconv.Atoi(pageStr)
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

		links, err := s.queries.UserFeed(r.Context(), database.UserFeedParams{
			UserID:   userID,
			Username: username,
			Limit:    25,
			Offset:   int32((page - 1) * 25),
		})
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		userpage.Page(userpage.Props{User: user, Links: links}).Render(r.Context(), w)
	}
}
