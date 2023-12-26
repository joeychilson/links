package server

import (
	"context"
	"net/http"

	"github.com/joeychilson/lixy/database"
	"github.com/joeychilson/lixy/pkg/sessions"
	"github.com/joeychilson/lixy/pkg/users"
)

func (s *Server) FetchCurrentUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cookie, err := s.sessions.Get(r)
		if err == nil {
			userID, err := s.queries.GetUserIDFromToken(ctx, database.GetUserIDFromTokenParams{
				Token:   cookie,
				Context: sessions.CookieName,
			})
			if err == nil {
				userRow, err := s.queries.GetUserByID(ctx, userID)
				if err == nil {
					ctx = context.WithValue(ctx, users.ContextKey, &users.User{
						ID:       userRow.ID,
						Email:    userRow.Email,
						Username: userRow.Username,
					})
				}
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) UserFromContext(ctx context.Context) *users.User {
	user, _ := ctx.Value(users.ContextKey).(*users.User)
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
