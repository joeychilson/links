package server

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/joeychilson/inquire/pages/home"
	"github.com/joeychilson/inquire/static"
)

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", static.Handler()))

	// Home page
	r.Get("/", templ.Handler(home.Page()).ServeHTTP)

	// Login page
	r.Route("/login", func(r chi.Router) {
		r.Get("/", s.handleLoginPage)
	})

	// Sign up page
	r.Route("/signup", func(r chi.Router) {
		r.Get("/", s.handleSignUpPage)
		r.Post("/", s.handleSignUp)
	})
	return r
}
