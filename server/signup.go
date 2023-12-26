package server

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

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
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm-password")

		emailExists, err := s.queries.EmailExists(r.Context(), email)
		if err != nil {
			log.Printf("Error checking email exists: %v\n", err)
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
			log.Printf("Error hashing password: %v\n", err)
			signup.Page(signup.PageProps{Error: ErrorInternalServer}).Render(r.Context(), w)
			return
		}

		userID, err := s.queries.CreateUser(r.Context(), database.CreateUserParams{
			Email:    email,
			Username: username,
			Password: string(hashedPassword),
		})
		if err != nil {
			log.Printf("Error creating user: %v\n", err)
			signup.Page(signup.PageProps{Error: ErrorInternalServer}).Render(r.Context(), w)
			return
		}

		err = s.sessionManager.Set(w, r, userID)
		if err != nil {
			log.Printf("Error setting session: %v\n", err)
			signup.Page(signup.PageProps{Error: ErrorInternalServer}).Render(r.Context(), w)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
