package accounts

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/BradiLugembaCS/payflow-api/internal/auth"
)

// Handler contains the database connection needed
// by our account-related routes.
type Handler struct {
	db *sql.DB
}

// NewHandler creates a new account handler.
func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

// CreateAccount handles POST /accounts.
//
// It creates an account for the currently logged-in user.
func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the authenticated user ID from the JWT middleware.
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(
			w,
			`{"error":"could not identify authenticated user"}`,
			http.StatusInternalServerError,
		)
		return
	}

	var account Account

	// Create a new account with a starting balance of 0.
	//
	// PostgreSQL automatically generates the account ID
	// and created_at timestamp.
	err := h.db.QueryRow(
		`
		INSERT INTO accounts (user_id)
		VALUES ($1)
		RETURNING id, user_id, balance, created_at
		`,
		userID,
	).Scan(
		&account.ID,
		&account.UserID,
		&account.Balance,
		&account.CreatedAt,
	)

	if err != nil {
		// Because user_id is UNIQUE, trying to create a second
		// account for the same user will fail.
		http.Error(
			w,
			`{"error":"could not create account"}`,
			http.StatusConflict,
		)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(account)
}

// GetMyAccount handles GET /accounts/me.
//
// It returns the account belonging to the logged-in user.
func (h *Handler) GetMyAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the user ID from the validated JWT.
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(
			w,
			`{"error":"could not identify authenticated user"}`,
			http.StatusInternalServerError,
		)
		return
	}

	var account Account

	// Look up the account that belongs to this user.
	err := h.db.QueryRow(
		`
		SELECT id, user_id, balance, created_at
		FROM accounts
		WHERE user_id = $1
		`,
		userID,
	).Scan(
		&account.ID,
		&account.UserID,
		&account.Balance,
		&account.CreatedAt,
	)

	if err == sql.ErrNoRows {
		http.Error(
			w,
			`{"error":"account not found"}`,
			http.StatusNotFound,
		)
		return
	}

	if err != nil {
		http.Error(
			w,
			`{"error":"could not retrieve account"}`,
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(account)
}
