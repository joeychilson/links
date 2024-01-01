package server

import (
	"net/http"

	"github.com/go-chi/httplog/v2"
	"github.com/joeychilson/links/components/feed"
	"github.com/joeychilson/links/db"
)

func (s *Server) PopularLinks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		linkRows, err := s.queries.PopularLinks(ctx, db.PopularLinksParams{
			UserID: user.ID,
			Offset: 0,
			Limit:  100,
		})
		if err != nil {
			oplog.Error("failed to get popular link feed", "error", err)
			feed.LinkFeed(feed.LinkFeedProps{User: user}).Render(ctx, w)
			return
		}

		oplog.Info("popular link feed", "count", len(linkRows))
		feed.LinkFeed(feed.LinkFeedProps{User: user, FeedType: feed.Popular, LinkRows: linkRows}).Render(ctx, w)
	}
}

func (s *Server) LatestLinks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		linkRows, err := s.queries.LatestLinks(ctx, db.LatestLinksParams{
			UserID: user.ID,
			Offset: 0,
			Limit:  100,
		})
		if err != nil {
			oplog.Error("failed to get popular link feed", "error", err)
			feed.LinkFeed(feed.LinkFeedProps{User: user}).Render(ctx, w)
			return
		}

		oplog.Info("latest link feed", "count", len(linkRows))
		feed.LinkFeed(feed.LinkFeedProps{User: user, FeedType: feed.Latest, LinkRows: linkRows}).Render(ctx, w)
	}
}

func (s *Server) ControversialLinks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		linkRows, err := s.queries.ControversialLinks(ctx, db.ControversialLinksParams{
			UserID: user.ID,
			Offset: 0,
			Limit:  100,
		})
		if err != nil {
			oplog.Error("failed to get popular link feed", "error", err)
			feed.LinkFeed(feed.LinkFeedProps{User: user}).Render(ctx, w)
			return
		}

		oplog.Info("controversial link feed", "count", len(linkRows))
		feed.LinkFeed(feed.LinkFeedProps{User: user, FeedType: feed.Controversial, LinkRows: linkRows}).Render(ctx, w)
	}
}
