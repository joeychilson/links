package server

import (
	"net/http"

	"github.com/go-chi/httplog/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/joeychilson/links/pages/login"
)

func (s *Server) LoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login.Page(login.Props{}).Render(r.Context(), w)
	}
}

func (s *Server) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)

		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validate email and password
		if email == "" || password == "" {
			login.Page(login.Props{Error: "Email and password are required"}).Render(ctx, w)
			return
		}

		// Attempt to log in
		user, err := s.queries.UserByEmail(ctx, email)
		if err != nil {
			oplog.Error("failed to get user by email", "error", err)
			login.Page(login.Props{Error: "Invalid email or password"}).Render(ctx, w)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			oplog.Error("failed to compare password", "error", err)
			login.Page(login.Props{Error: "Invalid email or password"}).Render(ctx, w)
			return
		}

		// Set session
		err = s.sessionManager.Set(w, r, user.ID)
		if err != nil {
			oplog.Error("failed to set session", "error", err)
			login.Page(login.Props{Error: ErrorInternalServer}).Render(ctx, w)
			return
		}

		oplog.Info("user logged in", "user_id", user.ID.String())
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
