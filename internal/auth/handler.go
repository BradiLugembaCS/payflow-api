package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Handler contains anything our authentication routes need.
//
// Right now, authentication needs access to the database,
// so we store our database connection here.
type Handler struct {
	db *sql.DB
}

// NewHandler creates a new authentication handler.
//
// We pass the database connection into it so registration
// and login can read and write users.
func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

// registerRequest represents the JSON we expect when
// somebody tries to create an account.
//
// Example:
//
//	{
//	    "email": "test@example.com",
//	    "password": "password123"
//	}
type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// loginRequest represents the JSON we expect
// when a user tries to log in.
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register handles POST /register requests.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {

	// Tell the client that our response will be JSON.
	w.Header().Set("Content-Type", "application/json")

	var req registerRequest

	// Convert the JSON request body into our Go struct.
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	// Remove accidental spaces from the email address.
	req.Email = strings.TrimSpace(req.Email)

	// Make sure an email address was provided.
	if req.Email == "" {
		http.Error(w, `{"error":"email is required"}`, http.StatusBadRequest)
		return
	}

	// Require a reasonably sized password.
	//
	// Later we can make our password rules stronger.
	if len(req.Password) < 8 {
		http.Error(
			w,
			`{"error":"password must be at least 8 characters"}`,
			http.StatusBadRequest,
		)
		return
	}

	// Hash the password using bcrypt.
	//
	// We NEVER save the original password in the database.
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		http.Error(
			w,
			`{"error":"could not process password"}`,
			http.StatusInternalServerError,
		)
		return
	}

	// Insert the new user into PostgreSQL.
	//
	// RETURNING id lets PostgreSQL give us the ID
	// that it generated for the new user.
	var userID int64

	err = h.db.QueryRow(
		`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id
		`,
		req.Email,
		string(passwordHash),
	).Scan(&userID)

	if err != nil {
		// For now we'll return a general registration error.
		//
		// Later we'll improve this so duplicate emails
		// return a more specific response.
		http.Error(
			w,
			`{"error":"could not create user"}`,
			http.StatusConflict,
		)
		return
	}

	// Registration succeeded.
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "user created successfully",
		"user_id": userID,
	})
}

// Login handles POST /login requests.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	// Tell the client that our response will be JSON.
	w.Header().Set("Content-Type", "application/json")

	var req loginRequest

	// Read the JSON request body.
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(
			w,
			`{"error":"invalid JSON"}`,
			http.StatusBadRequest,
		)
		return
	}

	// Remove accidental spaces from the email.
	req.Email = strings.TrimSpace(req.Email)

	// Basic validation.
	if req.Email == "" || req.Password == "" {
		http.Error(
			w,
			`{"error":"email and password are required"}`,
			http.StatusBadRequest,
		)
		return
	}

	// These variables will hold the user information
	// we retrieve from PostgreSQL.
	var userID int64
	var passwordHash string

	// Look up the user by email.
	//
	// We only need their ID and password hash for login.
	err = h.db.QueryRow(
		`
		SELECT id, password_hash
		FROM users
		WHERE email = $1
		`,
		req.Email,
	).Scan(&userID, &passwordHash)

	// If no user exists with that email,
	// return a generic login error.
	//
	// We don't say "email not found" because that can
	// reveal whether an account exists.
	if err == sql.ErrNoRows {
		http.Error(
			w,
			`{"error":"invalid email or password"}`,
			http.StatusUnauthorized,
		)
		return
	}

	// Handle unexpected database errors.
	if err != nil {
		http.Error(
			w,
			`{"error":"could not process login"}`,
			http.StatusInternalServerError,
		)
		return
	}

	// Compare the password the user entered
	// with the bcrypt hash stored in the database.
	err = bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(req.Password),
	)

	// If the password does not match, reject the login.
	if err != nil {
		http.Error(
			w,
			`{"error":"invalid email or password"}`,
			http.StatusUnauthorized,
		)
		return
	}

	// If we reach this point,
	// the email and password are correct.
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "login successful",
		"user_id": userID,
	})
}
