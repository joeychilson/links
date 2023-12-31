package server

import (
	"net/http"

	"github.com/go-chi/httplog/v2"
	"github.com/joeychilson/links/db"
	"github.com/joeychilson/links/pages/feed"
)

func (s *Server) PopularFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		linkRows, err := s.queries.PopularFeed(ctx, db.PopularFeedParams{
			UserID: user.ID,
			Offset: 0,
			Limit:  100,
		})
		if err != nil {
			oplog.Error("failed to get popular link feed", "error", err)
			feed.PopularFeed(feed.PopularFeedProps{User: user}).Render(ctx, w)
			return
		}

		oplog.Info("got popular link feed", "count", len(linkRows))
		feed.PopularFeed(feed.PopularFeedProps{User: user, LinkRows: linkRows}).Render(ctx, w)
	}
}

func (s *Server) LatestFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		linkRows, err := s.queries.LatestFeed(ctx, db.LatestFeedParams{
			UserID: user.ID,
			Offset: 0,
			Limit:  100,
		})
		if err != nil {
			oplog.Error("failed to get latest link feed", "error", err)
			feed.LatestFeed(feed.LatestFeedProps{User: user}).Render(ctx, w)
			return
		}

		oplog.Info("got latest link feed", "count", len(linkRows))
		feed.LatestFeed(feed.LatestFeedProps{User: user, LinkRows: linkRows}).Render(ctx, w)
	}
}

func (s *Server) ControversialFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		linkRows, err := s.queries.ControversialFeed(ctx, db.ControversialFeedParams{
			UserID: user.ID,
			Offset: 0,
			Limit:  100,
		})
		if err != nil {
			oplog.Error("failed to get controversial link feed", "error", err)
			feed.ControversialFeed(feed.ControversialFeedProps{User: user}).Render(ctx, w)
			return
		}

		oplog.Info("got latest link feed", "count", len(linkRows))
		feed.ControversialFeed(feed.ControversialFeedProps{User: user, LinkRows: linkRows}).Render(ctx, w)
	}
}
