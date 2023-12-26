package server

import (
	"net/http"
	"strconv"

	"github.com/joeychilson/lixy/database"
	"github.com/joeychilson/lixy/templates/pages/feed"
	"github.com/joeychilson/lixy/templates/pages/new"
)

func (s *Server) FeedPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := s.UserFromContext(r.Context())

		pageStr := r.URL.Query().Get("page")
		if pageStr == "" {
			pageStr = "1"
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			feed.Page(feed.Props{User: user}).Render(r.Context(), w)
			return
		}

		articleFeedRows, err := s.queries.ArticleFeed(r.Context(), database.ArticleFeedParams{
			Limit:  25,
			Offset: int32((page - 1) * 25),
		})
		if err != nil {
			feed.Page(feed.Props{User: user}).Render(r.Context(), w)
			return
		}

		feed.Page(feed.Props{User: user, ArticleFeedRows: articleFeedRows}).Render(r.Context(), w)
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
		user := s.UserFromContext(r.Context())

		title := r.FormValue("title")
		link := r.FormValue("link")

		if title == "" || link == "" {
			new.Page(new.PageProps{User: user}).Render(r.Context(), w)
			return
		}

		_, err := s.queries.CreateArticle(r.Context(), database.CreateArticleParams{
			UserID: user.ID,
			Title:  title,
			Link:   link,
		})
		if err != nil {
			new.Page(new.PageProps{User: user}).Render(r.Context(), w)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
