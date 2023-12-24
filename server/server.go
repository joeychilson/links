package server

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/joeychilson/flixmetrics/pages/home"
	"github.com/joeychilson/flixmetrics/pages/login"
	"github.com/joeychilson/flixmetrics/pages/signup"
)

type Server struct{}

func New() *Server {
	return &Server{}
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.Get("/", templ.Handler(home.Page()).ServeHTTP)

	r.Route("/login", func(r chi.Router) {
		r.Get("/", templ.Handler(login.Page()).ServeHTTP)
	})

	r.Route("/signup", func(r chi.Router) {
		r.Get("/", templ.Handler(signup.Page()).ServeHTTP)
	})
	return r
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.Router())
}
