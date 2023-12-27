package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/httplog/v2"

	"github.com/joeychilson/links/pkg/session"
)

func (s *Server) UserFromSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := s.sessionManager.GetUser(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		httplog.LogEntrySetField(ctx, "user_id", slog.StringValue(user.ID.String()))
		next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, session.SessionKey, user)))
	})
}

func (s *Server) UserFromContext(ctx context.Context) *session.User {
	user, _ := ctx.Value(session.SessionKey).(*session.User)
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
