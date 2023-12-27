package server

import (
	"net/http"

	"github.com/go-chi/httplog/v2"

	"github.com/joeychilson/links/database"
	"github.com/joeychilson/links/pages/me"
)

func (s *Server) MePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		user := s.UserFromContext(r.Context())

		feed, err := s.queries.LikedFeed(r.Context(), database.LikedFeedParams{
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
		me.Page(me.Props{User: user, Feed: feed}).Render(r.Context(), w)
	}
}
