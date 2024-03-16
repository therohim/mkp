CREATE TABLE IF NOT EXISTS transactions(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL,
    user_id UUID NOT NULL,
    seat_id UUID NOT NULL,
    service_fee FLOAT DEFAULT 0,
    total FLOAT DEFAULT 0,
    grand_total FLOAT DEFAULT 0,
    status INT DEFAULT 0,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (seat_id) REFERENCES ticket_seats(id)
);