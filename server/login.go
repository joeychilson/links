package server

import (
	"log"
	"net/http"

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
			log.Printf("Error getting user by email: %v\n", err)
			login.Page(login.Props{Error: "Invalid email or password"}).Render(r.Context(), w)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			log.Printf("Error comparing password: %v\n", err)
			login.Page(login.Props{Error: "Invalid email or password"}).Render(r.Context(), w)
			return
		}

		// Set session
		err = s.sessions.Set(w, r, user.ID)
		if err != nil {
			log.Printf("Error setting session: %v\n", err)
			login.Page(login.Props{Error: ErrorInternalServer}).Render(r.Context(), w)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
