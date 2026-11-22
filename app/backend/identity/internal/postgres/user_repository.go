
package postgres

import (
	"context"
	"errors"
	"fmt"

	"cardhub/backend/identity/internal/service"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(
	ctx context.Context,
	email, username, passwordHash string,
) (service.User, error) {
	var u service.User

	err := r.pool.QueryRow(ctx, `
		INSERT INTO users (email, username, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id::text, email, username, password_hash,
		          is_active, email_verified_at,
		          created_at, updated_at
	`, email, username, passwordHash).Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.PasswordHash,
		&u.IsActive,
		&u.EmailVerifiedAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return service.User{}, service.ErrUserExists
		}
		return service.User{}, fmt.Errorf("insert user: %w", err)
	}

	return u, nil
}

func (r *UserRepository) GetByID(
	ctx context.Context,
	id string,
) (service.User, error) {
	var u service.User

	err := r.pool.QueryRow(ctx, `
		SELECT id::text, email, username, password_hash,
		       is_active, email_verified_at,
		       created_at, updated_at
		FROM users
		WHERE id = $1::uuid
	`, id).Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.PasswordHash,
		&u.IsActive,
		&u.EmailVerifiedAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return service.User{}, service.ErrNotFound
	}
	if err != nil {
		return service.User{}, fmt.Errorf("select user by id: %w", err)
	}

	return u, nil
}

func (r *UserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (service.User, error) {
	var u service.User

	err := r.pool.QueryRow(ctx, `
		SELECT id::text, email, username, password_hash,
		       is_active, email_verified_at,
		       created_at, updated_at
		FROM users
		WHERE lower(email) = lower($1)
	`, email).Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.PasswordHash,
		&u.IsActive,
		&u.EmailVerifiedAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return service.User{}, service.ErrNotFound
	}
	if err != nil {
		return service.User{}, fmt.Errorf("select user by email: %w", err)
	}

	return u, nil
}
