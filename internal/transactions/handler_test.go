package transactions

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestCreateTransactionRequiresIdempotencyKey checks that payment requests are rejected when no Idempotency-Key header is supplied.
func TestCreateTransactionRequiresIdempotencyKey(t *testing.T) {
	// We don't need a database connection for this test because the request should fail before reaching the database.
	handler := &Handler{}

	// Create a fake payment request.
	req := httptest.NewRequest(
		http.MethodPost,
		"/transactions",
		strings.NewReader(`{
			"receiver_account_id": 3,
			"amount": 1000
		}`),
	)

	req.Header.Set("Content-Type", "application/json")

	// Capture the response from the handler.
	recorder := httptest.NewRecorder()

	// Run the transaction handler.
	handler.CreateTransaction(recorder, req)

	// A missing Idempotency-Key should return 400 Bad Request.
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}
}
