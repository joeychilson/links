package server

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"

	"github.com/joeychilson/links/db"
	notfound "github.com/joeychilson/links/pages/not_found"
	"github.com/joeychilson/links/pkg/session"
	"github.com/joeychilson/links/static"
)

// ErrorInternalServer is the error message to display when something goes wrong on the server
const ErrorInternalServer = "Sorry, something went wrong. Please try again later."

// Server represents the server of the application
type Server struct {
	logger         *httplog.Logger
	queries        *db.Queries
	sessionManager *session.Manager
}

// New returns a new server
func New(logger *httplog.Logger, queries *db.Queries, sessionManager *session.Manager) *Server {
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

	// Feeds
	r.Get("/", s.PopularFeed())
	r.Get("/latest", s.LatestFeed())
	r.Get("/controversial", s.ControversialFeed())

	// Link
	r.Route("/{slug}", func(r chi.Router) {
		r.Get("/", s.LinkPage())
		r.Route("/like", func(r chi.Router) {
			r.Use(s.RequireUser)
			r.Get("/", s.Like())
		})
		r.Route("/unlike", func(r chi.Router) {
			r.Use(s.RequireUser)
			r.Get("/", s.Unlike())
		})
	})

	// Create link
	r.Route("/create", func(r chi.Router) {
		r.Use(s.RequireUser)
		r.Get("/", s.CreateLinkPage())
		r.Post("/", s.CreateLink())
	})

	// Login
	r.Route("/login", func(r chi.Router) {
		r.Use(s.RedirectIfLoggedIn)
		r.Get("/", s.LogInPage())
		r.Post("/", s.LogIn())
	})

	// Logout
	r.Route("/logout", func(r chi.Router) {
		r.Use(s.RequireUser)
		r.Post("/", s.LogOut())
	})

	// Signup
	r.Route("/signup", func(r chi.Router) {
		r.Use(s.RedirectIfLoggedIn)
		r.Get("/", s.SignUpPage())
		r.Post("/", s.SignUp())
	})

	// Settings
	r.Route("/settings", func(r chi.Router) {
		r.Use(s.RequireUser)
		r.Get("/", s.SettingsPage())
	})

	// Not Found
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		user := s.UserFromContext(r.Context())
		notfound.Page(user).Render(r.Context(), w)
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

// RetargetPage is a helper function that makes retargeting to the whole page easier
func (s *Server) RetargetPage(ctx context.Context, w http.ResponseWriter, page templ.Component) {
	w.Header().Set("HX-Retarget", "#page")
	page.Render(ctx, w)
}
