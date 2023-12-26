package server

import (
	"context"
	"net/http"

	"github.com/joeychilson/lixy/database"
	"github.com/joeychilson/lixy/models"
)

type contextKey string

const userKey contextKey = "user"

func (s *Server) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cookie, err := r.Cookie("session")
		if err == nil {
			userID, err := s.queries.GetUserIDFromToken(ctx, database.GetUserIDFromTokenParams{
				Token:   cookie.Value,
				Context: "session",
			})
			if err == nil {
				userRow, err := s.queries.GetUserByID(ctx, userID)
				if err == nil {
					ctx = context.WithValue(ctx, userKey, &models.User{
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
