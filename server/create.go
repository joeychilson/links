package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/httplog/v2"
	"github.com/rs/xid"

	"github.com/joeychilson/links/db"
	"github.com/joeychilson/links/pages/create"
	"github.com/joeychilson/links/pkg/validate"
)

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
