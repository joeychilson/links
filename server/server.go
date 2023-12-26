package server

import (
	"net/http"

	"github.com/joeychilson/lixy/database"
	"github.com/joeychilson/lixy/pkg/sessions"
)

type Server struct {
	queries  *database.Queries
	sessions *sessions.Manager
}

func New(queries *database.Queries, sessions *sessions.Manager) *Server {
	return &Server{
		queries:  queries,
		sessions: sessions,
	}
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.Router())
}
