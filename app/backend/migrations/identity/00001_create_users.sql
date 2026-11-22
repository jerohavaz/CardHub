-- +goose Up

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    email TEXT NOT NULL,
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    email_verified_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_email_not_empty
        CHECK (length(trim(email)) > 0),

    CONSTRAINT users_username_format
        CHECK (username ~ '^[a-zA-Z0-9_]{3,30}$')
);

CREATE UNIQUE INDEX users_email_unique_idx
    ON users (lower(email));

CREATE UNIQUE INDEX users_username_unique_idx
    ON users (lower(username));

-- +goose Down

DROP TABLE users;