package server

import (
	"net/http"
	"strconv"

	"github.com/go-chi/httplog/v2"

	"github.com/joeychilson/links/components/linkfeed"
	"github.com/joeychilson/links/components/pagination"
	"github.com/joeychilson/links/db"
	"github.com/joeychilson/links/pages/feed"
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

		title := "Popular Feed"
		description := "Links that have been recently upvoted."

		linkCount, err := s.queries.CountLinks(ctx)
		if err != nil {
			oplog.Error("failed to get link count", "error", err)
			feed.Page(feed.Props{
				User:        user,
				Title:       title,
				Description: description,
				FeedType:    linkfeed.Popular,
			}).Render(ctx, w)
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
			feed.Page(feed.Props{
				User:        user,
				Title:       title,
				Description: description,
				FeedType:    linkfeed.Popular,
			}).Render(ctx, w)
			return
		}

		oplog.Info("popular link feed", "count", len(linkRows))
		props := feed.Props{
			User:        user,
			Title:       title,
			Description: description,
			FeedType:    linkfeed.Popular,
			LinkRows:    linkRows,
			Pagination: pagination.Props{
				CurrentPage: int64(page),
				TotalPages:  totalPages,
				Pages:       pages,
			},
		}
		feed.Page(props).Render(ctx, w)
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

		title := "Latest Feed"
		description := "Links that have been recently submitted."

		linkCount, err := s.queries.CountLinks(ctx)
		if err != nil {
			oplog.Error("failed to get link count", "error", err)
			feed.Page(feed.Props{
				User:        user,
				Title:       title,
				Description: description,
				FeedType:    linkfeed.Latest,
			}).Render(ctx, w)
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
			feed.Page(feed.Props{
				User:        user,
				Title:       title,
				Description: description,
				FeedType:    linkfeed.Latest,
			}).Render(ctx, w)
			return
		}

		oplog.Info("latest link feed", "count", len(linkRows))
		props := feed.Props{
			User:        user,
			Title:       title,
			Description: description,
			FeedType:    linkfeed.Latest,
			LinkRows:    linkRows,
			Pagination: pagination.Props{
				CurrentPage: int64(page),
				TotalPages:  totalPages,
				Pages:       pages,
			},
		}
		feed.Page(props).Render(ctx, w)
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

		title := "Controversial Feed"
		description := "Links that have been a lot of discussion."

		linkCount, err := s.queries.CountLinks(ctx)
		if err != nil {
			oplog.Error("failed to get link count", "error", err)
			feed.Page(feed.Props{
				User:        user,
				Title:       title,
				Description: description,
				FeedType:    linkfeed.Controversial,
			}).Render(ctx, w)
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
			feed.Page(feed.Props{
				User:        user,
				Title:       title,
				Description: description,
				FeedType:    linkfeed.Controversial,
			}).Render(ctx, w)
			return
		}

		oplog.Info("controversial link feed", "count", len(linkRows))
		props := feed.Props{
			User:        user,
			Title:       title,
			Description: description,
			FeedType:    linkfeed.Controversial,
			LinkRows:    linkRows,
			Pagination: pagination.Props{
				CurrentPage: int64(page),
				TotalPages:  totalPages,
				Pages:       pages,
			},
		}
		feed.Page(props).Render(ctx, w)
	}
}
