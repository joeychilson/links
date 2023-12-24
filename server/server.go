package server

import (
	"net/http"

	"github.com/joeychilson/flixmetrics/database"
)

type Server struct {
	queries *database.Queries
}

func New(queries *database.Queries) *Server {
	return &Server{
		queries: queries,
	}
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.Router())
}
