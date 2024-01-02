package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"

	"github.com/joeychilson/links/components/commentfeed"
	"github.com/joeychilson/links/components/link"
	"github.com/joeychilson/links/db"
	linkpage "github.com/joeychilson/links/pages/link"
)

func (s *Server) LinkPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)
		slug := chi.URLParam(r, "slug")

		linkRow, err := s.queries.LinkBySlug(ctx, db.LinkBySlugParams{
			UserID: user.ID,
			Slug:   slug,
		})
		if err != nil {
			oplog.Error("error getting link", err)
			s.Redirect(w, r, "/")
			return
		}

		commentRows, err := s.queries.PopularComments(ctx, db.PopularCommentsParams{
			Slug:   linkRow.Slug,
			UserID: user.ID,
		})
		if err != nil {
			oplog.Error("error getting comment feed", err)
			s.Redirect(w, r, "/")
			return
		}

		props := linkpage.Props{
			User:        user,
			LinkRow:     linkRow,
			FeedType:    commentfeed.Popular,
			CommentRows: commentRows,
			HasNextPage: linkRow.Comments > limit,
		}
		linkpage.Page(props).Render(ctx, w)
	}
}

func (s *Server) PopularComments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)
		slug := chi.URLParam(r, "slug")

		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			page = 1
		}

		commentRows, err := s.queries.PopularComments(ctx, db.PopularCommentsParams{
			Slug:   slug,
			UserID: user.ID,
		})
		if err != nil {
			oplog.Error("error getting comments", err)
			s.Redirect(w, r, fmt.Sprintf("/%s", slug))
			return
		}

		commentfeed.Nav(commentfeed.NavProps{LinkSlug: slug, Feed: commentfeed.Popular}).Render(ctx, w)
		feedProps := commentfeed.FeedProps{
			User:        user,
			FeedType:    commentfeed.Popular,
			CommentRows: commentRows,
		}
		commentfeed.Feed(feedProps).Render(ctx, w)
	}
}

func (s *Server) ControversialComments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)
		slug := chi.URLParam(r, "slug")

		commentRows, err := s.queries.ControversialComments(ctx, db.ControversialCommentsParams{
			Slug:   slug,
			UserID: user.ID,
		})
		if err != nil {
			oplog.Error("error getting comments", err)
			s.Redirect(w, r, fmt.Sprintf("/%s", slug))
			return
		}

		commentfeed.Nav(commentfeed.NavProps{LinkSlug: slug, Feed: commentfeed.Controversial}).Render(ctx, w)
		feedProps := commentfeed.FeedProps{
			User:        user,
			FeedType:    commentfeed.Controversial,
			CommentRows: commentRows,
		}
		commentfeed.Feed(feedProps).Render(ctx, w)
	}
}

func (s *Server) LatestComments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)
		slug := chi.URLParam(r, "slug")

		commentRows, err := s.queries.LatestComments(ctx, db.LatestCommentsParams{
			Slug:   slug,
			UserID: user.ID,
		})
		if err != nil {
			oplog.Error("error getting comments", err)
			s.Redirect(w, r, fmt.Sprintf("/%s", slug))
			return
		}

		var totalComments int64
		for _, comment := range commentRows {
			totalComments++
			totalComments += comment.Replies
		}

		commentfeed.Nav(commentfeed.NavProps{LinkSlug: slug, Feed: commentfeed.Latest}).Render(ctx, w)
		feedProps := commentfeed.FeedProps{
			User:        user,
			FeedType:    commentfeed.Latest,
			CommentRows: commentRows,
		}
		commentfeed.Feed(feedProps).Render(ctx, w)
	}
}

func (s *Server) Like() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)
		slug := chi.URLParam(r, "slug")

		linkID, err := s.queries.LinkIDBySlug(ctx, slug)
		if err != nil {
			oplog.Error("error getting link id", err)
			s.Redirect(w, r, "/")
			return
		}

		err = s.queries.CreateLike(ctx, db.CreateLikeParams{
			UserID: user.ID,
			LinkID: linkID,
		})
		if err != nil {
			oplog.Error("error creating like", err)
			s.Redirect(w, r, "/")
			return
		}

		linkRow, err := s.queries.LinkBySlug(ctx, db.LinkBySlugParams{
			UserID: user.ID,
			Slug:   slug,
		})
		if err != nil {
			oplog.Error("error getting link", err)
			s.Redirect(w, r, "/")
			return
		}

		oplog.Info("like created", "slug", slug)
		link.Component(link.Props{User: user, LinkRow: linkRow}).Render(ctx, w)
	}
}

func (s *Server) Unlike() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)
		slug := chi.URLParam(r, "slug")

		linkID, err := s.queries.LinkIDBySlug(ctx, slug)
		if err != nil {
			oplog.Error("error getting link id", err)
			s.Redirect(w, r, "/")
			return
		}

		err = s.queries.DeleteLike(ctx, db.DeleteLikeParams{
			UserID: user.ID,
			LinkID: linkID,
		})
		if err != nil {
			oplog.Error("error deleting like", err)
			s.Redirect(w, r, "/")
			return
		}

		linkRow, err := s.queries.LinkBySlug(ctx, db.LinkBySlugParams{
			UserID: user.ID,
			Slug:   slug,
		})
		if err != nil {
			oplog.Error("error getting link", err)
			s.Redirect(w, r, "/")
			return
		}

		oplog.Info("like deleted", "slug", slug)
		link.Component(link.Props{User: user, LinkRow: linkRow}).Render(ctx, w)
	}
}
