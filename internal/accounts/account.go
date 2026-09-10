package accounts

import "time"

// Account represents a payment account in our application.
//
// Each account belongs to one user and stores its balance
// as an integer number of cents.
type Account struct {
	ID int64 `json:"id"`

	UserID int64 `json:"user_id"`

	// Balance is stored in cents.
	//
	// Example:
	// 1050 means €10.50.
	Balance int64 `json:"balance"`

	CreatedAt time.Time `json:"created_at"`
}
