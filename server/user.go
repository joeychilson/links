package server

import (
	"net/http"
	"strconv"

	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"

	"github.com/joeychilson/links/database"
	userpage "github.com/joeychilson/links/pages/user"
)

func (s *Server) UserPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		username := r.URL.Query().Get("name")
		if username == "" {
			s.Redirect(w, r, "/")
			return
		}

		pageStr := r.URL.Query().Get("page")
		if pageStr == "" {
			pageStr = "1"
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			oplog.Error("failed to parse page number", "error", err)
			s.Redirect(w, r, "/")
			return
		}

		var userID uuid.UUID
		if user != nil {
			userID = user.ID
		} else {
			userID = uuid.Nil
		}

		feed, err := s.queries.UserFeedLinks(ctx, database.UserFeedLinksParams{
			UserID:   userID,
			Username: username,
			Limit:    25,
			Offset:   int32((page - 1) * 25),
		})
		if err != nil {
			oplog.Error("failed to get user links feed", "error", err)
			s.Redirect(w, r, "/")
			return
		}

		oplog.Info("user page loaded", "count", len(feed))
		userpage.Page(userpage.Props{User: user, Feed: feed}).Render(ctx, w)
	}
}
