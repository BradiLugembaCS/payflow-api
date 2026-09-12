-- Adding an idempotency key to each transaction.
--
-- The key lets us recognise repeated requests
-- and prevents the same payment from being created twice.
ALTER TABLE transactions
ADD COLUMN idempotency_key VARCHAR(255);

-- A key must be unique.
--
-- PostgreSQL will reject any attempt to create
-- another transaction using the same idempotency key.
CREATE UNIQUE INDEX IF NOT EXISTS
transactions_idempotency_key_unique
ON transactions(idempotency_key);