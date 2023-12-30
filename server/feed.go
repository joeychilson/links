package server

import (
	"net/http"

	"github.com/joeychilson/links/pages/feed"
)

func (s *Server) FeedPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := s.UserFromContext(ctx)

		feed.Page(&feed.Props{User: user}).Render(ctx, w)
	}
}
