CREATE TABLE IF NOT EXISTS ticket_studio(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bioskop VARCHAR(100) NULL,
    ticket_id UUID NOT NULL,
    studio VARCHAR(100) NULL,
    address TEXT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    FOREIGN KEY (ticket_id) REFERENCES tickets(id)
);