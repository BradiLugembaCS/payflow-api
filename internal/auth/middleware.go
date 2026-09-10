package auth

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// contextKey is a private type used for storing values
// inside the request context safely.
type contextKey string

const userIDKey contextKey = "userID"

// Middleware checks whether a request contains a valid JWT.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the Authorization header.
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(
				w,
				`{"error":"authorization token required"}`,
				http.StatusUnauthorized,
			)
			return
		}

		// The expected format is:
		//
		// Authorization: Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(
				w,
				`{"error":"invalid authorization header"}`,
				http.StatusUnauthorized,
			)
			return
		}

		tokenString := parts[1]

		secret := os.Getenv("JWT_SECRET")

		if secret == "" {
			http.Error(
				w,
				`{"error":"server authentication configuration missing"}`,
				http.StatusInternalServerError,
			)
			return
		}

		claims := &Claims{}

		// Parse the token and verify its signature.
		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			},
			jwt.WithValidMethods([]string{
				jwt.SigningMethodHS256.Alg(),
			}),
		)

		// Reject invalid, expired, or tampered tokens.
		if err != nil || !token.Valid {
			http.Error(
				w,
				`{"error":"invalid or expired token"}`,
				http.StatusUnauthorized,
			)
			return
		}

		// Store the user ID in the request context.
		//
		// The next handler can now access it.
		ctx := context.WithValue(
			r.Context(),
			userIDKey,
			claims.UserID,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserIDFromContext reads the authenticated user ID
// from the request context.
func UserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)

	return userID, ok
}
