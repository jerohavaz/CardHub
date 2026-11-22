
-- +goose Up

CREATE TABLE refresh_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    token_hash TEXT NOT NULL UNIQUE,

    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT refresh_sessions_expiry_check
        CHECK (expires_at > created_at)
);

CREATE INDEX refresh_sessions_user_id_idx
    ON refresh_sessions (user_id);


-- +goose Down

DROP TABLE IF EXISTS refresh_sessions;
