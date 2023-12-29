package server

import (
	"net/http"

	"github.com/go-chi/httplog/v2"

	"github.com/joeychilson/links/database"
	"github.com/joeychilson/links/pages/me"
)

func (s *Server) MePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		feed, err := s.queries.LikedFeed(ctx, database.LikedFeedParams{
			UserID: user.ID,
			Limit:  25,
			Offset: 0,
		})
		if err != nil {
			oplog.Error("failed to get liked feed", "error", err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		oplog.Info("me page loaded", "count", len(feed))
		me.Page(me.Props{User: user, Feed: feed}).Render(ctx, w)
	}
}
