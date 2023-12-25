package server

import (
	"fmt"
	"net/http"

	"github.com/joeychilson/inquire/pages/signup"
)

const (
	ErrorInternalServer = "internal server error"
	ErrorEmailExists    = "email is already in use"
	ErrorUsernameExists = "username is already in use"
)

func (s *Server) handleSignUpPage(w http.ResponseWriter, r *http.Request) {
	signup.Page().Render(r.Context(), w)
}

func (s *Server) handleSignUp(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm-password")

	fmt.Println(email, username, password, confirmPassword)
}

func (s *Server) handleEmailCheck(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	exists, err := s.queries.CheckEmailExists(r.Context(), email)
	if err != nil {
		signup.EmailInput(email, ErrorInternalServer).Render(r.Context(), w)
		return
	}

	if exists {
		signup.EmailInput(email, ErrorEmailExists).Render(r.Context(), w)
		return
	} else {
		signup.EmailInput(email, "").Render(r.Context(), w)
		return
	}
}

func (s *Server) handleUsernameCheck(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")

	exists, err := s.queries.CheckUsernameExists(r.Context(), username)
	if err != nil {
		signup.UsernameInput(username, ErrorInternalServer).Render(r.Context(), w)
		return
	}

	if exists {
		signup.UsernameInput(username, ErrorUsernameExists).Render(r.Context(), w)
		return
	} else {
		signup.UsernameInput(username, "").Render(r.Context(), w)
		return
	}
}
