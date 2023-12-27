package server

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-chi/httplog/v2"
	"github.com/joeychilson/links/database"
	"github.com/joeychilson/links/pages/signup"
)

var (
	ErrorInternalServer = "Sorry, something went wrong. Please try again later."
	ErrorEmailExists    = map[string]string{"email": "Sorry, this email is already in use"}
	ErrorUsernameExists = map[string]string{"username": "Sorry, this username is already in use"}
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
		oplog := httplog.LogEntry(r.Context())

		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm-password")

		emailExists, err := s.queries.EmailExists(r.Context(), email)
		if err != nil {
			oplog.Error("failed to check if email exists", "error", err)
			signup.Page(signup.PageProps{Error: ErrorInternalServer}).Render(r.Context(), w)
			return
		}

		if emailExists {
			formProps := signup.FormProps{
				Email:    email,
				Username: username,
				Error:    ErrorEmailExists,
			}
			signup.Page(signup.PageProps{FormProps: formProps}).Render(r.Context(), w)
			return
		}

		usernameExists, err := s.queries.UsernameExists(r.Context(), username)
		if err != nil {
			oplog.Error("failed to check if username exists", "error", err)
			signup.Page(signup.PageProps{Error: ErrorInternalServer}).Render(r.Context(), w)
			return
		}

		if usernameExists {
			formProps := signup.FormProps{
				Email:    email,
				Username: username,
				Error:    ErrorUsernameExists,
			}
			signup.Page(signup.PageProps{FormProps: formProps}).Render(r.Context(), w)
			return
		}

		if len(password) < 8 {
			formProps := signup.FormProps{
				Email:    email,
				Username: username,
				Error:    ErrorPasswordLength,
			}
			signup.Page(signup.PageProps{FormProps: formProps}).Render(r.Context(), w)
			return
		}

		if password != confirmPassword {
			formProps := signup.FormProps{
				Email:    email,
				Username: username,
				Error:    ErrorPasswordsMatch,
			}
			signup.Page(signup.PageProps{FormProps: formProps}).Render(r.Context(), w)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			oplog.Error("failed to hash password", "error", err)
			signup.Page(signup.PageProps{Error: ErrorInternalServer}).Render(r.Context(), w)
			return
		}

		userID, err := s.queries.CreateUser(r.Context(), database.CreateUserParams{
			Email:    email,
			Username: username,
			Password: string(hashedPassword),
		})
		if err != nil {
			oplog.Error("failed to create user", "error", err)
			signup.Page(signup.PageProps{Error: ErrorInternalServer}).Render(r.Context(), w)
			return
		}

		err = s.sessionManager.Set(w, r, userID)
		if err != nil {
			oplog.Error("failed to set session", "error", err)
			signup.Page(signup.PageProps{Error: ErrorInternalServer}).Render(r.Context(), w)
			return
		}

		oplog.Info("user signed up", "user_id", userID.String())
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
