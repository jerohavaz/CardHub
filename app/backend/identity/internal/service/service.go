package service

import (
	"context"
	"errors"
	"time"

	"cardhub/backend/identity/internal/token"
)

var (
	ErrNotFound           = errors.New("user not found")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidUsername    = errors.New("invalid username")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidID          = errors.New("invalid user ID")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountDisabled    = errors.New("account disabled")
	ErrInvalidSession     = errors.New("invalid or expired session")
)

type User struct {
	ID              string
	Email           string
	Username        string
	PasswordHash    string
	IsActive        bool
	EmailVerifiedAt *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// UserRepository handles user persistence.
type UserRepository interface {
	Create(
		ctx context.Context,
		email, username, passwordHash string,
	) (User, error)

	GetByID(ctx context.Context, id string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
}

// SessionRepository handles refresh-token persistence.
type SessionRepository interface {
	CreateSession(
		ctx context.Context,
		userID, tokenHash string,
		expiresAt time.Time,
	) error

	RotateSession(
		ctx context.Context,
		oldTokenHash, newTokenHash string,
		expiresAt time.Time,
	) (string, error)

	RevokeSession(
		ctx context.Context,
		tokenHash string,
	) error
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

type Service struct {
	users    UserRepository
	sessions SessionRepository
	signer   *token.Signer
}

func New(
	users UserRepository,
	sessions SessionRepository,
	signer *token.Signer,
) *Service {
	return &Service{
		users:    users,
		sessions: sessions,
		signer:   signer,
	}
}
