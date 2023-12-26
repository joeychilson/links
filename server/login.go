package server

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/joeychilson/lixy/database"
	"github.com/joeychilson/lixy/templates/pages/login"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) LoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login.Page(login.Props{}).Render(r.Context(), w)
	}
}

func (s *Server) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validate email and password
		if email == "" || password == "" {
			login.Page(login.Props{Error: "Email and password are required"}).Render(r.Context(), w)
			return
		}

		// Attempt to log in
		user, err := s.queries.GetUserByEmail(r.Context(), email)
		if err != nil {
			login.Page(login.Props{Error: "Invalid email or password"}).Render(r.Context(), w)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			login.Page(login.Props{Error: "Invalid email or password"}).Render(r.Context(), w)
			return
		}

		// Set session
		token, err := s.queries.CreateUserToken(r.Context(), database.CreateUserTokenParams{
			UserID:  user.ID,
			Token:   tokenGenerator(),
			Context: "session",
		})
		if err != nil {
			login.Page(login.Props{Error: "Sorry, something went wrong. Please try again later."}).Render(r.Context(), w)
			return
		}

		cookie := http.Cookie{
			Name:     "session",
			Value:    token,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func tokenGenerator() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
