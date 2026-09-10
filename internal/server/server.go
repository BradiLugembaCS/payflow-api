package server

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/BradiLugembaCS/payflow-api/internal/auth"
)

// Server represents our HTTP API.
//
// The router decides which function handles each URL.
type Server struct {
	router *http.ServeMux
}

// New creates and configures our HTTP server.
//
// We now pass in the database because some routes,
// such as registration, need database access.
func New(db *sql.DB) *Server {

	s := &Server{
		router: http.NewServeMux(),
	}

	// Create our authentication handler
	// and give it access to PostgreSQL.
	authHandler := auth.NewHandler(db)

	// Register all API routes.
	s.registerRoutes(authHandler)

	return s
}

// registerRoutes connects URLs to handler functions.
func (s *Server) registerRoutes(authHandler *auth.Handler) {

	// Existing health-check endpoint.
	s.router.HandleFunc("GET /health", s.handleHealth)

	// New registration endpoint.
	s.router.HandleFunc("POST /register", authHandler.Register)

	// Login endpoint.
	s.router.HandleFunc("POST /login", authHandler.Login)

	// Protected test endpoint.
	//
	// Only requests with a valid JWT can access this.
	s.router.Handle(
		"GET /me",
		auth.Middleware(http.HandlerFunc(s.handleMe)),
	)
}

// handleHealth lets us quickly confirm that the API is running.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

// Handler gives main.go access to our router.
func (s *Server) Handler() http.Handler {
	return s.router
}

// handleMe is a protected endpoint.
//
// It returns the ID of the currently authenticated user.
func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the user ID that the JWT middleware stored
	// inside the request context.
	userID, ok := auth.UserIDFromContext(r.Context())

	if !ok {
		http.Error(
			w,
			`{"error":"user not found in request context"}`,
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userID,
	})
}
