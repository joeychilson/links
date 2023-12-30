package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"

	"github.com/joeychilson/links/db"
	"github.com/joeychilson/links/static"
)

// Server represents the server of the application
type Server struct {
	logger  *httplog.Logger
	queries *db.Queries
}

// New returns a new server
func New(logger *httplog.Logger, queries *db.Queries) *Server {
	return &Server{
		logger:  logger,
		queries: queries,
	}
}

// Router returns the http.Handler for the server
// This is where we define all of our routes
func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(httplog.RequestLogger(s.logger))
	r.Use(middleware.Recoverer)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", static.Handler()))

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
