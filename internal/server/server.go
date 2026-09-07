package server

import (
	"encoding/json"
	"net/http"
)

type Server struct {
	router *http.ServeMux
}

func New() *Server {
	s := &Server{
		router: http.NewServeMux(),
	}

	s.registerRoutes()

	return s
}

func (s *Server) registerRoutes() {
	s.router.HandleFunc("GET /health", s.handleHealth)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (s *Server) Handler() http.Handler {
	return s.router
}