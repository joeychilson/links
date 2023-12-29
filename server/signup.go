package server

import (
	"net/http"

	"github.com/go-chi/httplog/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/joeychilson/links/database"
	"github.com/joeychilson/links/pages/signup"
)

var (
	ErrorInternalServer = "Sorry, something went wrong. Please try again later."
	ErrorEmailExists    = map[string]string{"email": "Sorry, this email is already in use"}
	ErrorUsernameExists = map[string]string{"username": "Sorry, this username is already in use"}
	ErrorUsernameLength = map[string]string{"username": "Username must be at least 4 characters"}
	ErrorPasswordLength = map[string]string{"password": "Password must be at least 8 characters"}
	ErrorPasswordSymbol = map[string]string{"password": "Password must contain at least one symbol"}
	ErrorPasswordsMatch = map[string]string{"confirm-password": "Passwords do not match"}
)

func (s *Server) SignUpPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		signup.Page(signup.PageProps{}).Render(r.Context(), w)
	}
}

func (s *Server) SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)

		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm-password")

		emailExists, err := s.queries.EmailExists(ctx, email)
		if err != nil {
			oplog.Error("failed to check if email exists", "error", err)
			signup.Page(signup.PageProps{Error: ErrorInternalServer}).Render(ctx, w)
			return
		}

		if emailExists {
			props := signup.FormProps{
				Email:    email,
				Username: username,
				Error:    ErrorEmailExists,
			}
			signup.Form(props).Render(ctx, w)
			return
		}

		usernameExists, err := s.queries.UsernameExists(ctx, username)
		if err != nil {
			oplog.Error("failed to check if username exists", "error", err)
			signup.Page(signup.PageProps{Error: ErrorInternalServer}).Render(ctx, w)
			return
		}

		if usernameExists {
			props := signup.FormProps{
				Email:    email,
				Username: username,
				Error:    ErrorUsernameExists,
			}
			signup.Form(props).Render(ctx, w)
			return
		}

		if len(username) < 4 {
			props := signup.FormProps{
				Email:    email,
				Username: username,
				Error:    ErrorUsernameLength,
			}
			signup.Form(props).Render(ctx, w)
			return
		}

		if len(password) < 8 {
			props := signup.FormProps{
				Email:    email,
				Username: username,
				Error:    ErrorPasswordLength,
			}
			signup.Form(props).Render(ctx, w)
			return
		}

		if password != confirmPassword {
			props := signup.FormProps{
				Email:    email,
				Username: username,
				Error:    ErrorPasswordsMatch,
			}
			signup.Form(props).Render(ctx, w)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			oplog.Error("failed to hash password", "error", err)
			signup.Page(signup.PageProps{Error: ErrorInternalServer}).Render(ctx, w)
			return
		}

		userID, err := s.queries.CreateUser(ctx, database.CreateUserParams{
			Email:    email,
			Username: username,
			Password: string(hashedPassword),
		})
		if err != nil {
			oplog.Error("failed to create user", "error", err)
			signup.Page(signup.PageProps{Error: ErrorInternalServer}).Render(ctx, w)
			return
		}

		err = s.sessionManager.Set(w, r, userID)
		if err != nil {
			oplog.Error("failed to set session", "error", err)
			signup.Page(signup.PageProps{Error: ErrorInternalServer}).Render(ctx, w)
			return
		}

		oplog.Info("user signed up", "user_id", userID.String())
		s.Redirect(w, "/")
	}
}
