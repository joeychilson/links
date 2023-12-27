package server

import (
	"net/http"
	"strconv"

	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"

	"github.com/joeychilson/links/database"
	"github.com/joeychilson/links/pages/feed"
	"github.com/joeychilson/links/pages/new"
)

func (s *Server) FeedPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		user := s.UserFromContext(r.Context())

		pageStr := r.URL.Query().Get("page")
		if pageStr == "" {
			pageStr = "1"
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			oplog.Error("failed to parse page number", "error", err)
			feed.Page(feed.Props{User: user}).Render(r.Context(), w)
			return
		}

		var userID uuid.UUID
		if user != nil {
			userID = user.ID
		} else {
			userID = uuid.Nil
		}

		feedRows, err := s.queries.LinkFeed(r.Context(), database.LinkFeedParams{
			UserID: userID,
			Limit:  25,
			Offset: int32((page - 1) * 25),
		})
		if err != nil {
			oplog.Error("failed to get link feed", "error", err)
			feed.Page(feed.Props{User: user}).Render(r.Context(), w)
			return
		}

		oplog.Info("feed page loaded", "count", len(feedRows))
		feed.Page(feed.Props{User: user, Feed: feedRows}).Render(r.Context(), w)
	}
}

func (s *Server) NewPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := s.UserFromContext(r.Context())
		new.Page(new.PageProps{User: user}).Render(r.Context(), w)
	}
}

func (s *Server) New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		user := s.UserFromContext(r.Context())

		title := r.FormValue("title")
		url := r.FormValue("url")

		if title == "" || url == "" {
			new.Page(new.PageProps{User: user}).Render(r.Context(), w)
			return
		}

		err := s.queries.CreateLink(r.Context(), database.CreateLinkParams{
			UserID: user.ID,
			Title:  title,
			Url:    url,
		})
		if err != nil {
			oplog.Error("failed to create link", "error", err)
			new.Page(new.PageProps{User: user}).Render(r.Context(), w)
			return
		}

		oplog.Info("user created link", "title", title, "url", url)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
