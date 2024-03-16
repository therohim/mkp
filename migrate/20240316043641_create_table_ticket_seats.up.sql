CREATE TABLE IF NOT EXISTS ticket_seats(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_id UUID NOT NULL,
    seat VARCHAR(50) DEFAULT NULL,
    is_buy BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    FOREIGN KEY (schedule_id) REFERENCES ticket_schedule(id)
);