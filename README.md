# PayFlow API

PayFlow API simulates the core backend of a simple payment platform. Users can register and log in, create a payment account, view their balance, and send money to other accounts through authenticated API requests. The application also includes safeguards that are important in payment systems, such as secure password hashing, JWT-based authentication, database transactions, row locking, and idempotency keys to help prevent duplicate or inconsistent payments.

I built this project to strengthen my backend engineering skills and learn Go while exploring concepts that are important in payment systems, such as authentication, database transactions, concurrency control, and idempotency.

## Features

- User registration and login
- Password hashing with bcrypt
- JWT authentication
- Protected API routes
- User payment accounts
- Account balances stored in cents
- Account-to-account transfers
- PostgreSQL database transactions
- Row locking with `FOR UPDATE`
- Idempotency keys to prevent duplicate payments
- Basic API validation
- Automated handler tests

## Tech Stack

- Go
- PostgreSQL
- pgx
- JWT
- bcrypt
- Go `net/http`
- Go `httptest`
- Git and GitHub

## Project Structure

```text
payflow-api/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── accounts/
│   ├── auth/
│   ├── database/
│   ├── server/
│   └── transactions/
│
├── migrations/
├── .env.example
├── go.mod
└── README.md