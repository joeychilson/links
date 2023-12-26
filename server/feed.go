package server

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/joeychilson/links/database"
	"github.com/joeychilson/links/templates/components/article"
	"github.com/joeychilson/links/templates/pages/feed"
	"github.com/joeychilson/links/templates/pages/new"
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

		err := s.queries.CreateArticle(r.Context(), database.CreateArticleParams{
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

func (s *Server) Like() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := s.UserFromContext(r.Context())
		articleID := r.URL.Query().Get("articleID")

		if articleID == "" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		articleUUID, err := uuid.Parse(articleID)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		err = s.queries.CreateLike(r.Context(), database.CreateLikeParams{
			UserID:    user.ID,
			ArticleID: articleUUID,
		})
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		likeCount, err := s.queries.CountLikes(r.Context(), articleUUID)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		article.VoteCount(articleUUID, likeCount).Render(r.Context(), w)
	}
}
