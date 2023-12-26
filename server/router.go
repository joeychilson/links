package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/joeychilson/lixy/static"
	"github.com/joeychilson/lixy/templates/pages/home"
)

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(s.FetchCurrentUser)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", static.Handler()))

	// Home page
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		user := s.UserFromContext(r.Context())
		home.Page(home.Props{User: user}).Render(r.Context(), w)
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
