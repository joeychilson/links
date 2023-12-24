package server

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/joeychilson/starter-templ/pages/home"
	"github.com/joeychilson/starter-templ/pages/login"
	"github.com/joeychilson/starter-templ/pages/signup"
)

type Server struct{}

func New() *Server {
	return &Server{}
}

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	mux.HandleFunc("/", templ.Handler(home.Page()).ServeHTTP)
	mux.HandleFunc("/login", templ.Handler(login.Page()).ServeHTTP)
	mux.HandleFunc("/signup", templ.Handler(signup.Page()).ServeHTTP)
	return mux
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.Router())
}
