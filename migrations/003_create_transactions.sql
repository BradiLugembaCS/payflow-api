-- Stores every transfer between two accounts.
CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,

    -- Account sending the money.
    sender_account_id BIGINT NOT NULL,

    -- Account receiving the money.
    receiver_account_id BIGINT NOT NULL,

    -- Amount is stored in cents.
    amount BIGINT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Both account IDs must exist.
    CONSTRAINT fk_transactions_sender
        FOREIGN KEY (sender_account_id)
        REFERENCES accounts(id),

    CONSTRAINT fk_transactions_receiver
        FOREIGN KEY (receiver_account_id)
        REFERENCES accounts(id),

    -- A transaction must move a positive amount.
    CONSTRAINT transactions_amount_positive
        CHECK (amount > 0),

    -- Prevent sending money to the same account.
    CONSTRAINT transactions_different_accounts
        CHECK (sender_account_id <> receiver_account_id)
);