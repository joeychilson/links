package server

import (
	"context"
	"net/http"

	"github.com/joeychilson/inquire/database"
	"github.com/joeychilson/inquire/models"
)

type contextKey string

const userKey contextKey = "user"

func (s *Server) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx context.Context

		// Attempt to retrieve the session cookie
		cookie, err := r.Cookie("session")
		if err == nil {
			// If the cookie is found, try to get the user
			userID, err := s.queries.GetUserIDFromToken(r.Context(), database.GetUserIDFromTokenParams{
				Token:   cookie.Value,
				Context: "session",
			})
			if err == nil {
				user, err := s.queries.GetUserByID(r.Context(), userID)
				if err == nil {
					// If user is found, add user to the context
					ctx = context.WithValue(r.Context(), userKey, models.User{
						ID:       user.ID,
						Email:    user.Email,
						Username: user.Username,
					})
				}
			}
		}

		// If ctx is not set, use the original request context
		if ctx == nil {
			ctx = context.WithValue(r.Context(), userKey, models.User{})
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
