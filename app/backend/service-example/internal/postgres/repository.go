// Package postgres is the only example-service package that knows SQL.
package postgres

import (
	"context"
	"errors"
	"fmt"

	"cardhub/backend/service-example/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct{ pool *pgxpool.Pool }

func New(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

func (r *Repository) Create(ctx context.Context, name string) (service.Example, error) {
	var id string
	err := r.pool.QueryRow(ctx, `INSERT INTO examples (name) VALUES ($1) RETURNING id::text`, name).Scan(&id)
	if err != nil {
		return service.Example{}, fmt.Errorf("insert example: %w", err)
	}
	return service.Example{ID: id, Name: name}, nil
}

func (r *Repository) Get(ctx context.Context, id string) (service.Example, error) {
	var e service.Example
	err := r.pool.QueryRow(ctx, `SELECT id::text, name FROM examples WHERE id = $1::uuid`, id).Scan(&e.ID, &e.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return service.Example{}, service.ErrNotFound
	}
	if err != nil {
		return service.Example{}, fmt.Errorf("select example: %w", err)
	}
	return e, nil
}
