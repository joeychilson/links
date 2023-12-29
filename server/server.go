package server

import (
	"net/http"

	"github.com/go-chi/httplog/v2"
	"github.com/joeychilson/links/database"
	"github.com/joeychilson/links/pkg/session"
)

type Server struct {
	logger         *httplog.Logger
	queries        *database.Queries
	sessionManager *session.Manager
}

func New(logger *httplog.Logger, queries *database.Queries, sessionManager *session.Manager) *Server {
	return &Server{
		logger:         logger,
		queries:        queries,
		sessionManager: sessionManager,
	}
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.Router())
}

func (s *Server) Redirect(w http.ResponseWriter, path string) {
	w.Header().Set("HX-Redirect", path)
	w.WriteHeader(http.StatusOK)
}
