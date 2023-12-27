package server

import (
	"net/http"

	"github.com/joeychilson/links/database"
	"github.com/joeychilson/links/pkg/session"
)

type Server struct {
	queries        *database.Queries
	sessionManager *session.Manager
}

func New(queries *database.Queries, sessionManager *session.Manager) *Server {
	return &Server{
		queries:        queries,
		sessionManager: sessionManager,
	}
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.Router())
}
