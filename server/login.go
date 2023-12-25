package server

import (
	"log"
	"net/http"

	"github.com/joeychilson/inquire/database"
	"github.com/joeychilson/inquire/models"
	"github.com/joeychilson/inquire/pages/login"
	"github.com/joeychilson/inquire/pages/signup"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userKey).(models.User)
	if user.ID != 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	login.Page(login.Props{}).Render(r.Context(), w)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
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
		Token:   s.tokenGenerator(),
		Context: "session",
	})
	if err != nil {
		log.Printf("Error creating user token: %v", err)
		signup.Page(signup.PageProps{Error: ErrorInternalServer}).Render(r.Context(), w)
		return
	}

	cookie := http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusFound)
}
