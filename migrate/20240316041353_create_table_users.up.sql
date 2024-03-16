CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NULL,
    email VARCHAR (50) NULL,
    phone VARCHAR (50) NULL,
    password VARCHAR (255) NULL,
    birthday DATE NULL,
    photo VARCHAR (255) NULL,
    status INT DEFAULT 1,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);