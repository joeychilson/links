package server

import (
	"net/http"

	"github.com/joeychilson/inquire/pages/signup"
)

const (
	ErrorInternalServer = "Sorry, something went wrong. Please try again later."
	ErrorEmailExists    = "Sorry, this email is already in use"
	ErrorUsernameExists = "Sorry, This username is already in use"
)

func (s *Server) handleSignUpPage(w http.ResponseWriter, r *http.Request) {
	signup.Page().Render(r.Context(), w)
}

func (s *Server) handleSignUp(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	username := r.FormValue("username")

	if email == "testing@test.com" {
		signup.EmailField(signup.EmailFieldProps{Email: email, Error: ErrorEmailExists}).Render(r.Context(), w)
		return
	}

	if username == "testing" {
		signup.UsernameField(signup.UsernameFieldProps{Username: username, Error: ErrorUsernameExists}).Render(r.Context(), w)
		return
	}
}

func (s *Server) handleCheckEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	exists, err := s.queries.CheckEmailExists(r.Context(), email)
	if err != nil {
		// TODO: Log error, and make it render an alert instead of input error.
		signup.EmailField(signup.EmailFieldProps{Email: email, Error: ErrorInternalServer}).Render(r.Context(), w)
		return
	}

	if exists {
		signup.EmailField(signup.EmailFieldProps{Email: email, Error: ErrorEmailExists}).Render(r.Context(), w)
		return
	} else {
		signup.EmailField(signup.EmailFieldProps{Email: email}).Render(r.Context(), w)
		return
	}
}

func (s *Server) handleCheckUsername(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")

	exists, err := s.queries.CheckUsernameExists(r.Context(), username)
	if err != nil {
		// TODO: Log error, and make it render an alert instead of input error.
		signup.UsernameField(signup.UsernameFieldProps{Username: username, Error: ErrorInternalServer}).Render(r.Context(), w)
		return
	}

	if exists {
		signup.UsernameField(signup.UsernameFieldProps{Username: username, Error: ErrorUsernameExists}).Render(r.Context(), w)
		return
	} else {
		signup.UsernameField(signup.UsernameFieldProps{Username: username}).Render(r.Context(), w)
		return
	}
}
