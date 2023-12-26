package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/joeychilson/lixy/models"
	"github.com/joeychilson/lixy/pages/home"
	"github.com/joeychilson/lixy/static"
)

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(s.Authorization)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", static.Handler()))

	// Home page
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Context().Value(userKey).(*models.User)
		home.Page(home.Props{User: user}).Render(r.Context(), w)
	})

	// Account page
	r.Route("/account", func(r chi.Router) {
		r.Get("/", s.handleAccountPage)
	})

	// Login page
	r.Route("/login", func(r chi.Router) {
		r.Get("/", s.handleLoginPage)
		r.Post("/", s.handleLogin)
	})

	// Logout
	r.Post("/logout", s.handleLogout)

	// Sign up page
	r.Route("/signup", func(r chi.Router) {
		r.Get("/", s.handleSignUpPage)
		r.Post("/", s.handleSignUp)
	})
	return r
}
