package server

import (
	"net/http"
	"strconv"

	"github.com/go-chi/httplog/v2"

	"github.com/joeychilson/links/components/linkfeed"
	"github.com/joeychilson/links/db"
	"github.com/joeychilson/links/pages/feed"
)

const limit = 25

func (s *Server) FeedPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		linkRows, err := s.queries.PopularLinks(ctx, db.PopularLinksParams{
			UserID: user.ID,
			Offset: 0,
			Limit:  limit,
		})
		if err != nil {
			oplog.Error("failed to get popular link feed", "error", err)
			s.Redirect(w, r, "/")
			return
		}

		oplog.Info("popular link feed", "count", len(linkRows))
		props := feed.Props{
			User:        user,
			FeedType:    linkfeed.Popular,
			LinkRows:    linkRows,
			HasNextPage: len(linkRows) == limit,
		}
		feed.Page(props).Render(ctx, w)
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
		linkfeed.Nav(linkfeed.NavProps{Feed: linkfeed.Popular}).Render(ctx, w)
		props := linkfeed.FeedProps{
			User:        user,
			LinkRows:    linkRows,
			FeedType:    linkfeed.Popular,
			NextPage:    page + 1,
			HasNextPage: len(linkRows) == limit,
		}
		linkfeed.Feed(props).Render(ctx, w)
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
		linkfeed.Nav(linkfeed.NavProps{Feed: linkfeed.Latest}).Render(ctx, w)
		props := linkfeed.FeedProps{
			User:        user,
			LinkRows:    linkRows,
			FeedType:    linkfeed.Latest,
			NextPage:    page + 1,
			HasNextPage: len(linkRows) == limit,
		}
		linkfeed.Feed(props).Render(ctx, w)
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
		linkfeed.Nav(linkfeed.NavProps{Feed: linkfeed.Controversial}).Render(ctx, w)
		props := linkfeed.FeedProps{
			User:        user,
			LinkRows:    linkRows,
			FeedType:    linkfeed.Controversial,
			NextPage:    page + 1,
			HasNextPage: len(linkRows) == limit,
		}
		linkfeed.Feed(props).Render(ctx, w)
	}
}
