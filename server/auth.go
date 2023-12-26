package server

import (
	"context"
	"net/http"

	"github.com/joeychilson/lixy/pkg/session"
)

func (s *Server) UserFromSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := s.sessionManager.GetUser(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx = context.WithValue(ctx, session.ContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) UserFromContext(ctx context.Context) *session.User {
	user, _ := ctx.Value(session.ContextKey).(*session.User)
	return user
}

func (s *Server) RedirectIfLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.UserFromContext(r.Context())
		if user != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.UserFromContext(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
