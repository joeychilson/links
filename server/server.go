package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"

	"github.com/joeychilson/links/database"
	"github.com/joeychilson/links/pages/errors"
	"github.com/joeychilson/links/pkg/session"
	"github.com/joeychilson/links/static"
)

// Server represents the server of the application
type Server struct {
	logger         *httplog.Logger
	queries        *database.Queries
	sessionManager *session.Manager
}

// New returns a new server
func New(logger *httplog.Logger, queries *database.Queries, sessionManager *session.Manager) *Server {
	return &Server{
		logger:         logger,
		queries:        queries,
		sessionManager: sessionManager,
	}
}

// Router returns the http.Handler for the server
// This is where we define all of our routes
func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(httplog.RequestLogger(s.logger))
	r.Use(middleware.Recoverer)
	r.Use(s.UserFromSession)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", static.Handler()))

	// Feed page
	r.Route("/", func(r chi.Router) {
		r.Get("/", s.FeedPage())
	})

	// New post page
	r.Route("/new", func(r chi.Router) {
		r.Use(s.RequireUser)
		r.Get("/", s.NewPage())
		r.Post("/", s.New())
	})

	// Link
	r.Route("/link", func(r chi.Router) {
		r.Get("/", s.Link())
	})

	// Comment
	r.Route("/comment", func(r chi.Router) {
		r.Use(s.RequireUser)
		r.Post("/", s.Comment())
	})

	// Reply
	r.Route("/reply", func(r chi.Router) {
		r.Use(s.RequireUser)
		r.Get("/", s.CommentReply())
	})

	// Vote
	r.Route("/vote", func(r chi.Router) {
		r.Use(s.RequireUser)
		r.Get("/", s.Vote())
	})

	// User
	r.Route("/user", func(r chi.Router) {
		r.Get("/", s.UserPage())
	})

	// Me page
	r.Route("/me", func(r chi.Router) {
		r.Use(s.RequireUser)
		r.Get("/", s.MePage())
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

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		user := s.UserFromContext(r.Context())
		errors.NotFound(user).Render(r.Context(), w)
	})
	return r
}

// Redirect is a helper function that makes redirects easier with HX-Request
func (s *Server) Redirect(w http.ResponseWriter, r *http.Request, path string) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", path)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Redirect(w, r, path, http.StatusFound)
	}
}
