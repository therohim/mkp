CREATE TABLE IF NOT EXISTS tickets(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NULL,
    duration INT DEFAULT 0,
    genre TEXT DEFAULT NULL,
    cover varchar(100) DEFAULT NULL,
    description TEXT DEFAULT NULL,
    rating FLOAT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);