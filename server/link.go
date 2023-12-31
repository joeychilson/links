package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"

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

		var userID uuid.UUID
		if user != nil {
			userID = user.ID
		} else {
			userID = uuid.Nil
		}

		dbLink, err := s.queries.LinkBySlug(ctx, db.LinkBySlugParams{
			Column1: userID,
			Slug:    slug,
		})
		if err != nil {
			oplog.Error("error getting link", err)
			s.Redirect(w, r, "/")
			return
		}

		commentRows, err := s.queries.CommentFeed(ctx, db.CommentFeedParams{
			Slug:   dbLink.Slug,
			UserID: userID,
			Offset: 0,
			Limit:  100,
		})
		if err != nil {
			oplog.Error("error getting comment feed", err)
			s.Redirect(w, r, "/")
			return
		}

		linkpage.Page(&linkpage.Props{User: user, Link: dbLink, CommentRows: commentRows}).Render(ctx, w)
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
			Column1: user.ID,
			Slug:    slug,
		})
		if err != nil {
			oplog.Error("error getting link", err)
			s.Redirect(w, r, "/")
			return
		}

		oplog.Info("like created", "slug", slug)
		link.Component(&link.Props{User: user, LinkRow: link.LinkRow(linkRow)}).Render(ctx, w)
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
			Column1: user.ID,
			Slug:    slug,
		})
		if err != nil {
			oplog.Error("error getting link", err)
			s.Redirect(w, r, "/")
			return
		}

		oplog.Info("like deleted", "slug", slug)
		link.Component(&link.Props{User: user, LinkRow: link.LinkRow(linkRow)}).Render(ctx, w)
	}
}
