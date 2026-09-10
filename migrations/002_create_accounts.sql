-- Each user can have a payment account.
-- The balance is stored in cents to avoid floating-point errors.
CREATE TABLE IF NOT EXISTS accounts (
    id BIGSERIAL PRIMARY KEY,

    -- Links this account to a user.
    user_id BIGINT NOT NULL UNIQUE,

    -- Store money as an integer number of cents.
    -- Example: €12.50 is stored as 1250.
    balance BIGINT NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- If the user is deleted, their account is deleted too.
    CONSTRAINT fk_accounts_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    -- Prevent negative balances at the database level.
    CONSTRAINT accounts_balance_non_negative
        CHECK (balance >= 0)
);