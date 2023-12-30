package server

import (
	"net/http"

	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"
	"github.com/joeychilson/links/db"
	"github.com/joeychilson/links/pages/feed"
)

func (s *Server) FeedPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		var userID uuid.UUID
		if user != nil {
			userID = user.ID
		} else {
			userID = uuid.Nil
		}

		feedRows, err := s.queries.LinkFeed(ctx, db.LinkFeedParams{
			Column1: userID,
			Offset:  0,
			Limit:   100,
		})
		if err != nil {
			oplog.Error("failed to get link feed", "error", err)
			feed.Page(&feed.Props{User: user}).Render(ctx, w)
			return
		}

		oplog.Info("got link feed", "count", len(feedRows))
		feed.Page(&feed.Props{User: user, FeedRows: feedRows}).Render(ctx, w)
	}
}
