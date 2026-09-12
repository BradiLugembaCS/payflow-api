package transactions

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/BradiLugembaCS/payflow-api/internal/auth"
)

// Handler contains the database connection used by transaction routes.
type Handler struct {
	db *sql.DB
}

// NewHandler creates a transaction handler.
func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

// createTransactionRequest represents the JSON body
// used when sending money.
type createTransactionRequest struct {
	ReceiverAccountID int64 `json:"receiver_account_id"`
	Amount            int64 `json:"amount"`
}

// CreateTransaction handles POST /transactions.
func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Identify the logged-in user from the JWT.
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(
			w,
			`{"error":"could not identify authenticated user"}`,
			http.StatusInternalServerError,
		)
		return
	}

	var req createTransactionRequest

	// Read the JSON request body.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(
			w,
			`{"error":"invalid JSON"}`,
			http.StatusBadRequest,
		)
		return
	}

	// The amount must be greater than zero.
	if req.Amount <= 0 {
		http.Error(
			w,
			`{"error":"amount must be greater than zero"}`,
			http.StatusBadRequest,
		)
		return
	}

	// Start a database transaction.
	//
	// This means all money changes either succeed together
	// or fail together.
	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(
			w,
			`{"error":"could not start transaction"}`,
			http.StatusInternalServerError,
		)
		return
	}

	// If something fails later, rollback will undo changes.
	defer tx.Rollback()

	var senderAccountID int64
	var senderBalance int64

	// Find the logged-in user's account.
	//
	// FOR UPDATE locks this account row while the transfer is happening.
	// This helps protect against concurrent balance changes.
	err = tx.QueryRowContext(
		r.Context(),
		`
		SELECT id, balance
		FROM accounts
		WHERE user_id = $1
		FOR UPDATE
		`,
		userID,
	).Scan(
		&senderAccountID,
		&senderBalance,
	)

	if err == sql.ErrNoRows {
		http.Error(
			w,
			`{"error":"sender account not found"}`,
			http.StatusNotFound,
		)
		return
	}

	if err != nil {
		http.Error(
			w,
			`{"error":"could not retrieve sender account"}`,
			http.StatusInternalServerError,
		)
		return
	}

	// Prevent sending money to yourself.
	if senderAccountID == req.ReceiverAccountID {
		http.Error(
			w,
			`{"error":"cannot send money to the same account"}`,
			http.StatusBadRequest,
		)
		return
	}

	// Check the sender has enough money.
	if senderBalance < req.Amount {
		http.Error(
			w,
			`{"error":"insufficient balance"}`,
			http.StatusBadRequest,
		)
		return
	}

	// Check that the receiver exists and lock that row too.
	var receiverAccountID int64

	err = tx.QueryRowContext(
		r.Context(),
		`
		SELECT id
		FROM accounts
		WHERE id = $1
		FOR UPDATE
		`,
		req.ReceiverAccountID,
	).Scan(&receiverAccountID)

	if err == sql.ErrNoRows {
		http.Error(
			w,
			`{"error":"receiver account not found"}`,
			http.StatusNotFound,
		)
		return
	}

	if err != nil {
		http.Error(
			w,
			`{"error":"could not retrieve receiver account"}`,
			http.StatusInternalServerError,
		)
		return
	}

	// Deduct money from the sender.
	_, err = tx.ExecContext(
		r.Context(),
		`
		UPDATE accounts
		SET balance = balance - $1
		WHERE id = $2
		`,
		req.Amount,
		senderAccountID,
	)

	if err != nil {
		http.Error(
			w,
			`{"error":"could not update sender balance"}`,
			http.StatusInternalServerError,
		)
		return
	}

	// Add money to the receiver.
	_, err = tx.ExecContext(
		r.Context(),
		`
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
		`,
		req.Amount,
		receiverAccountID,
	)

	if err != nil {
		http.Error(
			w,
			`{"error":"could not update receiver balance"}`,
			http.StatusInternalServerError,
		)
		return
	}

	var transaction Transaction

	// Record the transfer.
	err = tx.QueryRowContext(
		r.Context(),
		`
		INSERT INTO transactions (
			sender_account_id,
			receiver_account_id,
			amount
		)
		VALUES ($1, $2, $3)
		RETURNING id, sender_account_id, receiver_account_id, amount, created_at
		`,
		senderAccountID,
		receiverAccountID,
		req.Amount,
	).Scan(
		&transaction.ID,
		&transaction.SenderAccountID,
		&transaction.ReceiverAccountID,
		&transaction.Amount,
		&transaction.CreatedAt,
	)

	if err != nil {
		http.Error(
			w,
			`{"error":"could not record transaction"}`,
			http.StatusInternalServerError,
		)
		return
	}

	// Commit makes all database changes permanent.
	if err := tx.Commit(); err != nil {
		http.Error(
			w,
			`{"error":"could not complete transaction"}`,
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(transaction)
}
