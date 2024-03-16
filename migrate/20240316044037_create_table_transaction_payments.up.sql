CREATE TABLE IF NOT EXISTS transaction_payments(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL,
    channel VARCHAR(100) NOT NULL,
    amount FLOAT DEFAULT 0,
    status INT DEFAULT 0,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    FOREIGN KEY (transaction_id) REFERENCES transactions(id)
);