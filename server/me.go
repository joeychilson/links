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

		feedType := r.URL.Query().Get("feed")

		var (
			feed []database.FeedRow
			err  error
		)
		if feedType == "voted" {
			feed, err = s.queries.UserFeedVoted(ctx, database.UserFeedVotedParams{
				UserID: user.ID,
				Limit:  25,
				Offset: 0,
			})
			if err != nil {
				oplog.Error("failed to get user voted feed", "error", err)
				s.Redirect(w, r, "/")
				return
			}
		} else {
			feed, err = s.queries.UserFeedLinks(ctx, database.UserFeedLinksParams{
				UserID:   user.ID,
				Username: user.Username,
				Limit:    25,
				Offset:   0,
			})
			if err != nil {
				oplog.Error("failed to get user feed", "error", err)
				s.Redirect(w, r, "/")
				return
			}
		}

		oplog.Info("me page loaded", "feed_type", feedType, "count", len(feed))
		me.Page(me.Props{User: user, FeedType: feedType, Feed: feed}).Render(ctx, w)
	}
}
