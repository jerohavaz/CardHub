package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"cardhub/backend/identity/internal/password"
)

const refreshTokenBytes = 32

var (
	usernamePattern = regexp.MustCompile(
		`^[a-zA-Z0-9_]{3,30}$`,
	)

	uuidPattern = regexp.MustCompile(
		`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
	)
)

// Register validates input, hashes the password
// and creates a new user.
func (s *Service) Register(
	ctx context.Context,
	email, username, plainPassword string,
) (User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	username = strings.TrimSpace(username)

	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email ||
		len(email) > 254 {
		return User{}, ErrInvalidEmail
	}

	if !usernamePattern.MatchString(username) {
		return User{}, ErrInvalidUsername
	}

	if len(plainPassword) < 12 ||
		len(plainPassword) > 128 {
		return User{}, ErrInvalidPassword
	}

	hash, err := password.Hash(plainPassword)
	if err != nil {
		return User{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.users.Create(
		ctx, email, username, hash,
	)
	if err != nil {
		return User{}, err
	}

	user.PasswordHash = ""
	return user, nil
}

// Login verifies credentials and issues an access JWT
// together with a persistent refresh session.
func (s *Service) Login(
	ctx context.Context,
	email, plainPassword string,
) (AuthTokens, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return AuthTokens{}, ErrInvalidCredentials
		}
		return AuthTokens{}, err
	}

	valid, err := password.Verify(
		plainPassword, user.PasswordHash,
	)
	if err != nil {
		return AuthTokens{}, fmt.Errorf(
			"verify password: %w", err,
		)
	}
	if !valid {
		return AuthTokens{}, ErrInvalidCredentials
	}

	if !user.IsActive {
		return AuthTokens{}, ErrAccountDisabled
	}

	accessToken, err := s.signer.Issue(user.ID)
	if err != nil {
		return AuthTokens{}, fmt.Errorf(
			"issue access token: %w", err,
		)
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return AuthTokens{}, err
	}

	err = s.sessions.CreateSession(
		ctx,
		user.ID,
		hashRefreshToken(refreshToken),
		time.Now().UTC().Add(s.refreshTTL),
	)
	if err != nil {
		return AuthTokens{}, fmt.Errorf(
			"create session: %w", err,
		)
	}

	return AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GetMe retrieves the authenticated user.
// userID must come from trusted authentication context.
func (s *Service) GetMe(
	ctx context.Context,
	userID string,
) (User, error) {
	if !uuidPattern.MatchString(userID) {
		return User{}, ErrInvalidID
	}

	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return User{}, err
	}

	if !user.IsActive {
		return User{}, ErrAccountDisabled
	}

	user.PasswordHash = ""
	return user, nil
}

// RefreshToken rotates the refresh token
// and issues a new access JWT.
func (s *Service) RefreshToken(
	ctx context.Context,
	oldRefreshToken string,
) (AuthTokens, error) {
	if !validRefreshToken(oldRefreshToken) {
		return AuthTokens{}, ErrInvalidSession
	}

	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		return AuthTokens{}, err
	}

	userID, err := s.sessions.RotateSession(
		ctx,
		hashRefreshToken(oldRefreshToken),
		hashRefreshToken(newRefreshToken),
		time.Now().UTC().Add(s.refreshTTL),
	)
	if err != nil {
		return AuthTokens{}, err
	}

	accessToken, err := s.signer.Issue(userID)
	if err != nil {
		return AuthTokens{}, fmt.Errorf(
			"issue access token: %w", err,
		)
	}

	return AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// Logout revokes the supplied refresh session.
// The operation is idempotent.
func (s *Service) Logout(
	ctx context.Context,
	refreshToken string,
) error {
	if !validRefreshToken(refreshToken) {
		return ErrInvalidSession
	}

	return s.sessions.RevokeSession(
		ctx,
		hashRefreshToken(refreshToken),
	)
}

// generateRefreshToken creates a random 256-bit token.
func generateRefreshToken() (string, error) {
	b := make([]byte, refreshTokenBytes)

	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf(
			"generate refresh token: %w", err,
		)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

// hashRefreshToken hashes a high-entropy token
// before database storage.
func hashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// validRefreshToken checks the expected token format.
func validRefreshToken(token string) bool {
	b, err := base64.RawURLEncoding.DecodeString(token)
	return err == nil && len(b) == refreshTokenBytes &&
		base64.RawURLEncoding.EncodeToString(b) == token
}
