package server

import (
	"net/http"

	"github.com/joeychilson/starter-templ/pages/home"
)

type Server struct{}

func New() *Server {
	return &Server{}
}

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		home.Page().Render(r.Context(), w)
	})

	return mux
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.Router())
}
