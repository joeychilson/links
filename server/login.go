package server

import (
	"net/http"

	"github.com/go-chi/httplog/v2"
	"github.com/joeychilson/links/pages/login"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) LogInPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login.Page(&login.Props{}).Render(r.Context(), w)
	}
}

func (s *Server) LogIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)

		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			props := &login.Props{
				Error: "Please enter your email and password.",
			}
			login.Page(props).Render(ctx, w)
			return
		}

		userIDPasswordRow, err := s.queries.UserIDAndPasswordByEmail(ctx, email)
		if err != nil {
			oplog.Error("failed to get user by email", "error", err)
			props := &login.Props{
				Error: "Please provide a valid email and password.",
			}
			login.Page(props).Render(ctx, w)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(userIDPasswordRow.Password), []byte(password))
		if err != nil {
			if err == bcrypt.ErrMismatchedHashAndPassword {
				props := &login.Props{
					Error: "Please provide a valid email and password.",
				}
				login.Page(props).Render(ctx, w)
				return
			}
			oplog.Error("failed to compare password", "error", err)
			return
		}

		err = s.sessionManager.Set(w, r, userIDPasswordRow.ID)
		if err != nil {
			oplog.Error("failed to set session", "error", err)
			props := &login.Props{
				Error: ErrorInternalServer,
			}
			login.Page(props).Render(ctx, w)
			return
		}

		oplog.Info("user logged in", "user_id", userIDPasswordRow.ID.String())
		s.Redirect(w, r, "/")
	}
}
