package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/rs/xid"

	"github.com/joeychilson/links/components/link"
	"github.com/joeychilson/links/db"
	"github.com/joeychilson/links/pages/create"
	linkpage "github.com/joeychilson/links/pages/link"
	"github.com/joeychilson/links/pkg/validate"
)

func (s *Server) LinkPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		slug := chi.URLParam(r, "slug")
		if slug == "" {
			s.Redirect(w, r, "/")
			return
		}

		dbLink, err := s.queries.LinkBySlug(ctx, slug)
		if err != nil {
			oplog.Error("error getting link", err)
			s.Redirect(w, r, "/")
			return
		}

		linkpage.Page(&linkpage.Props{User: user, Link: dbLink}).Render(ctx, w)
	}
}

func (s *Server) CreateLinkPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := s.UserFromContext(ctx)
		create.Page(&create.Props{User: user, FormProps: &create.FormProps{}}).Render(ctx, w)
	}
}

func (s *Server) CreateLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		title := r.FormValue("title")
		link := r.FormValue("link")

		validationError := validate.Title(title)
		if validationError != nil {
			props := &create.FormProps{
				Title: title,
				Link:  link,
				Error: validationError,
			}
			create.Form(props).Render(ctx, w)
			return
		}

		validationError = validate.Link(link)
		if validationError != nil {
			props := &create.FormProps{
				Title: title,
				Link:  link,
				Error: validationError,
			}
			create.Form(props).Render(ctx, w)
			return
		}

		slug, err := s.queries.CreateLink(ctx, db.CreateLinkParams{
			UserID: user.ID,
			Title:  title,
			Url:    link,
			Slug:   xid.New().String(),
		})
		if err != nil {
			oplog.Error("error creating link", err)
			props := &create.Props{
				Error:     ErrorInternalServer,
				FormProps: &create.FormProps{},
			}
			s.RetargetPage(ctx, w, create.Page(props))
			return
		}

		oplog.Info("link created", "slug", slug)
		s.Redirect(w, r, fmt.Sprintf("/%s", slug))
	}
}

func (s *Server) Like() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		slug := chi.URLParam(r, "slug")
		if slug == "" {
			s.Redirect(w, r, "/")
			return
		}

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

		likesAndLiked, err := s.queries.LinkLikesAndLiked(ctx, db.LinkLikesAndLikedParams{
			Column1: linkID,
			Column2: user.ID,
		})
		if err != nil {
			oplog.Error("error getting likes and liked", err)
			s.Redirect(w, r, "/")
			return
		}

		oplog.Info("like created", "slug", slug)
		link.LikeButton(slug, likesAndLiked.Liked).Render(ctx, w)
		link.Likes(slug, likesAndLiked.Likes).Render(ctx, w)
	}
}

func (s *Server) Unlike() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		slug := chi.URLParam(r, "slug")
		if slug == "" {
			s.Redirect(w, r, "/")
			return
		}

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

		likesAndLiked, err := s.queries.LinkLikesAndLiked(ctx, db.LinkLikesAndLikedParams{
			Column1: linkID,
			Column2: user.ID,
		})
		if err != nil {
			oplog.Error("error getting likes and liked", err)
			s.Redirect(w, r, "/")
			return
		}

		oplog.Info("like deleted", "slug", slug)
		link.LikeButton(slug, likesAndLiked.Liked).Render(ctx, w)
		link.Likes(slug, likesAndLiked.Likes).Render(ctx, w)
	}
}
