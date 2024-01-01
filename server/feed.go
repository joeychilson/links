package server

import (
	"net/http"
	"strconv"

	"github.com/go-chi/httplog/v2"
	"github.com/joeychilson/links/components/feed"
	"github.com/joeychilson/links/db"
)

const (
	limit          = 25
	maxPagesToShow = 5
)

func (s *Server) Feed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			page = 1
		}

		linkRows, err := s.queries.PopularLinks(ctx, db.PopularLinksParams{
			UserID: user.ID,
			Offset: int32((page - 1) * limit),
			Limit:  int32(limit),
		})
		if err != nil {
			oplog.Error("failed to get popular link feed", "error", err)
			s.Redirect(w, r, "/")
			return
		}

		oplog.Info("popular link feed", "count", len(linkRows))
		props := feed.LinkFeedProps{
			User:        user,
			Title:       "Popular Feed",
			Description: "Links that have been recently upvoted.",
			FeedType:    feed.Popular,
			LinkRows:    linkRows,
			HasNextPage: len(linkRows) == limit,
		}
		feed.LinkFeed(props).Render(ctx, w)
	}
}

func (s *Server) PopularLinks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			page = 1
		}

		linkRows, err := s.queries.PopularLinks(ctx, db.PopularLinksParams{
			UserID: user.ID,
			Offset: int32((page - 1) * limit),
			Limit:  int32(limit),
		})
		if err != nil {
			oplog.Error("failed to get popular link feed", "error", err)
			s.Redirect(w, r, "/")
			return
		}

		oplog.Info("popular link feed", "count", len(linkRows))
		props := feed.FeedProps{
			User:        user,
			LinkRows:    linkRows,
			FeedType:    feed.Popular,
			NextPage:    page + 1,
			HasNextPage: len(linkRows) == limit,
		}
		feed.LinkFeedNav(feed.LinkFeedNavProps{Feed: feed.Popular}).Render(ctx, w)
		feed.Feed(props).Render(ctx, w)
	}
}

func (s *Server) LatestLinks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			page = 1
		}

		linkRows, err := s.queries.LatestLinks(ctx, db.LatestLinksParams{
			UserID: user.ID,
			Offset: int32((page - 1) * limit),
			Limit:  int32(limit),
		})
		if err != nil {
			oplog.Error("failed to get latest link feed", "error", err)
			s.Redirect(w, r, "/")
			return
		}

		oplog.Info("latest link feed", "count", len(linkRows))
		props := feed.FeedProps{
			User:        user,
			LinkRows:    linkRows,
			FeedType:    feed.Latest,
			NextPage:    page + 1,
			HasNextPage: len(linkRows) == limit,
		}
		feed.LinkFeedNav(feed.LinkFeedNavProps{Feed: feed.Latest}).Render(ctx, w)
		feed.Feed(props).Render(ctx, w)
	}
}

func (s *Server) ControversialLinks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			page = 1
		}

		linkRows, err := s.queries.ControversialLinks(ctx, db.ControversialLinksParams{
			UserID: user.ID,
			Offset: int32((page - 1) * limit),
			Limit:  int32(limit),
		})
		if err != nil {
			oplog.Error("failed to get controversial link feed", "error", err)
			s.Redirect(w, r, "/")
			return
		}

		oplog.Info("controversial link feed", "count", len(linkRows))
		props := feed.FeedProps{
			User:        user,
			LinkRows:    linkRows,
			FeedType:    feed.Controversial,
			NextPage:    page + 1,
			HasNextPage: len(linkRows) == limit,
		}
		feed.LinkFeedNav(feed.LinkFeedNavProps{Feed: feed.Controversial}).Render(ctx, w)
		feed.Feed(props).Render(ctx, w)
	}
}
