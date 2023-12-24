package server

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/joeychilson/flixmetrics/pages/home"
	"github.com/joeychilson/flixmetrics/pages/signup"
)

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.Get("/", templ.Handler(home.Page()).ServeHTTP)

	r.Route("/signup", func(r chi.Router) {
		r.Get("/", templ.Handler(signup.Page()).ServeHTTP)
		r.Post("/check-email", func(w http.ResponseWriter, r *http.Request) {
			email := r.FormValue("email")
			if email != "testing@test.com" {
				signup.EmailInput(false, email).Render(r.Context(), w)
			} else {
				signup.EmailInput(true, email).Render(r.Context(), w)
			}
		})
		r.Post("/check-username", func(w http.ResponseWriter, r *http.Request) {
			username := r.FormValue("username")
			if username != "testing" {
				signup.UsernameInput(false, username).Render(r.Context(), w)
			} else {
				signup.UsernameInput(true, username).Render(r.Context(), w)
			}
		})
	})
	return r
}
