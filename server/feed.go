package server

import (
	"net/http"
	"strconv"

	"github.com/joeychilson/lixy/database"
	"github.com/joeychilson/lixy/templates/pages/feed"
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
