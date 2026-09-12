package transactions

import "time"

// Transaction represents a money transfer between two accounts.
type Transaction struct {
	ID int64 `json:"id"`

	SenderAccountID int64 `json:"sender_account_id"`

	ReceiverAccountID int64 `json:"receiver_account_id"`

	// Amount is stored in cents.
	Amount int64 `json:"amount"`

	CreatedAt time.Time `json:"created_at"`
}
