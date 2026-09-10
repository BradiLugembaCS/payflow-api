package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the information we store inside the JWT.
//
// We only store information that is safe to include.
// Never put passwords or other secrets inside a JWT.
type Claims struct {
	UserID int64 `json:"user_id"`

	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT for a user.
func GenerateToken(userID int64) (string, error) {
	// Read the JWT signing secret from the environment.
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not set")
	}

	// Set an expiry time.
	//
	// For now, tokens will last for 1 hour.
	expiresAt := time.Now().Add(1 * time.Hour)

	// Create the information that will go inside the token.
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create a new token using HMAC SHA-256.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using our secret key.
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
