package server

import (
	"net/http"
	"strconv"

	"github.com/go-chi/httplog/v2"
	"github.com/joeychilson/links/components/feed"
	"github.com/joeychilson/links/components/pagination"
	"github.com/joeychilson/links/db"
)

const (
	limit          = 25
	maxPagesToShow = 5
)

func (s *Server) PopularLinks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			page = 1
		}

		linkCount, err := s.queries.CountLinks(ctx)
		if err != nil {
			oplog.Error("failed to get link count", "error", err)
			feed.LinkFeed(feed.LinkFeedProps{User: user}).Render(ctx, w)
			return
		}

		totalPages := (linkCount + limit - 1) / limit
		offset := (page - 1) * limit
		pages := pagination.Pages(int64(page), totalPages, maxPagesToShow)

		linkRows, err := s.queries.PopularLinks(ctx, db.PopularLinksParams{
			UserID: user.ID,
			Offset: int32(offset),
			Limit:  int32(limit),
		})
		if err != nil {
			oplog.Error("failed to get popular link feed", "error", err)
			feed.LinkFeed(feed.LinkFeedProps{User: user}).Render(ctx, w)
			return
		}

		oplog.Info("popular link feed", "count", len(linkRows))
		props := feed.LinkFeedProps{
			User:        user,
			Title:       "Popular Feed",
			Description: "Links that have been recently upvoted.",
			FeedType:    feed.Popular,
			LinkRows:    linkRows,
			Pagination: pagination.Props{
				CurrentPage: int64(page),
				TotalPages:  totalPages,
				Pages:       pages,
			},
		}
		feed.LinkFeed(props).Render(ctx, w)
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

		linkCount, err := s.queries.CountLinks(ctx)
		if err != nil {
			oplog.Error("failed to get link count", "error", err)
			feed.LinkFeed(feed.LinkFeedProps{User: user}).Render(ctx, w)
			return
		}

		totalPages := (linkCount + limit - 1) / limit
		offset := (page - 1) * limit
		pages := pagination.Pages(int64(page), totalPages, maxPagesToShow)

		linkRows, err := s.queries.LatestLinks(ctx, db.LatestLinksParams{
			UserID: user.ID,
			Offset: int32(offset),
			Limit:  int32(limit),
		})
		if err != nil {
			oplog.Error("failed to get popular link feed", "error", err)
			feed.LinkFeed(feed.LinkFeedProps{User: user}).Render(ctx, w)
			return
		}

		oplog.Info("latest link feed", "count", len(linkRows))
		props := feed.LinkFeedProps{
			User:        user,
			Title:       "Latest Feed",
			Description: "Links that have been recently submitted.",
			FeedType:    feed.Latest,
			LinkRows:    linkRows,
			Pagination: pagination.Props{
				CurrentPage: int64(page),
				TotalPages:  totalPages,
				Pages:       pages,
			},
		}
		feed.LinkFeed(props).Render(ctx, w)
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

		linkCount, err := s.queries.CountLinks(ctx)
		if err != nil {
			oplog.Error("failed to get link count", "error", err)
			feed.LinkFeed(feed.LinkFeedProps{User: user}).Render(ctx, w)
			return
		}

		totalPages := (linkCount + limit - 1) / limit
		offset := (page - 1) * limit
		pages := pagination.Pages(int64(page), totalPages, maxPagesToShow)

		linkRows, err := s.queries.ControversialLinks(ctx, db.ControversialLinksParams{
			UserID: user.ID,
			Offset: int32(offset),
			Limit:  int32(limit),
		})
		if err != nil {
			oplog.Error("failed to get popular link feed", "error", err)
			feed.LinkFeed(feed.LinkFeedProps{User: user}).Render(ctx, w)
			return
		}

		oplog.Info("controversial link feed", "count", len(linkRows))
		props := feed.LinkFeedProps{
			User:        user,
			Title:       "Controversial Feed",
			Description: "Links that have a lot of comments with a lot of disagreement.",
			FeedType:    feed.Controversial,
			LinkRows:    linkRows,
			Pagination: pagination.Props{
				CurrentPage: int64(page),
				TotalPages:  totalPages,
				Pages:       pages,
			},
		}
		feed.LinkFeed(props).Render(ctx, w)
	}
}
