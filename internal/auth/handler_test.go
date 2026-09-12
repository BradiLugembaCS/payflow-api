package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRegisterRejectsShortPassword checks that registration
// fails when the password is shorter than 8 characters.
func TestRegisterRejectsShortPassword(t *testing.T) {
	// We do not need a real database connection for this test
	// because the request should fail before reaching PostgreSQL.
	handler := &Handler{}

	// Create a fake HTTP request.
	req := httptest.NewRequest(
		http.MethodPost,
		"/register",
		strings.NewReader(`{
			"email":"test@example.com",
			"password":"short"
		}`),
	)

	req.Header.Set("Content-Type", "application/json")

	// Recorder captures the HTTP response.
	recorder := httptest.NewRecorder()

	// Run the handler.
	handler.Register(recorder, req)

	// We expect HTTP 400 Bad Request.
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}
}

// TestRegisterRejectsInvalidJSON checks that malformed JSON
// is rejected instead of being processed.
func TestRegisterRejectsInvalidJSON(t *testing.T) {
	handler := &Handler{}

	req := httptest.NewRequest(
		http.MethodPost,
		"/register",
		strings.NewReader(`{"email":`),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	handler.Register(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}
}
