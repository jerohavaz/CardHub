package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cardhub/backend/identity/internal/service"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	pool *pgxpool.Pool
}

func NewSessionRepository(pool *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{pool: pool}
}

func (r *SessionRepository) CreateSession(
	ctx context.Context,
	userID, tokenHash string,
	expiresAt time.Time,
) error {
	result, err := r.pool.Exec(ctx, `
		INSERT INTO refresh_sessions (
			user_id, token_hash, expires_at
		)
		SELECT id, $2, $3
		FROM users
		WHERE id = $1::uuid
		  AND is_active = TRUE
	`, userID, tokenHash, expiresAt)

	if err != nil {
		return fmt.Errorf("create refresh session: %w", err)
	}

	if result.RowsAffected() == 0 {
		return service.ErrInvalidSession
	}

	return nil
}

func (r *SessionRepository) RotateSession(
	ctx context.Context,
	oldTokenHash, newTokenHash string,
	expiresAt time.Time,
) (string, error) {
	var userID string

	err := r.pool.QueryRow(ctx, `
		UPDATE refresh_sessions AS s
		SET token_hash = $2,
		    expires_at = $3
		FROM users AS u
		WHERE s.token_hash = $1
		  AND s.user_id = u.id
		  AND u.is_active = TRUE
		  AND s.revoked_at IS NULL
		  AND s.expires_at > NOW()
		  AND $3 > NOW()
		RETURNING s.user_id::text
	`, oldTokenHash, newTokenHash, expiresAt).Scan(&userID)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", service.ErrInvalidSession
	}
	if err != nil {
		return "", fmt.Errorf("rotate session: %w", err)
	}

	return userID, nil
}

func (r *SessionRepository) RevokeSession(
	ctx context.Context,
	tokenHash string,
) error {
	_, err := r.pool.Exec(ctx, `
        UPDATE refresh_sessions
        SET revoked_at = NOW()
        WHERE token_hash = $1
          AND revoked_at IS NULL
    `, tokenHash)

	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}

	return nil
}
