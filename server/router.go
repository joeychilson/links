package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/joeychilson/lixy/static"
)

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(s.UserFromSession)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", static.Handler()))

	// Feed page
	r.Route("/", func(r chi.Router) {
		r.Get("/", s.FeedPage())
	})

	// Account page
	r.Route("/account", func(r chi.Router) {
		r.Use(s.RequireUser)
		r.Get("/", s.AccountPage())
	})

	// Login page
	r.Route("/login", func(r chi.Router) {
		r.Use(s.RedirectIfLoggedIn)
		r.Get("/", s.LoginPage())
		r.Post("/", s.Login())
	})

	// Logout
	r.Post("/logout", s.Logout())

	// Sign up page
	r.Route("/signup", func(r chi.Router) {
		r.Use(s.RedirectIfLoggedIn)
		r.Get("/", s.SignUpPage())
		r.Post("/", s.SignUp())
	})
	return r
}
