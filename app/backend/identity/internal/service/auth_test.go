package service

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"cardhub/backend/identity/internal/password"
	"cardhub/backend/identity/internal/token"
)

const (
	userID       = "123e4567-e89b-12d3-a456-426614174000"
	goodPassword = "SecurePassword123!"

	testIssuer     = "cardhub-identity"
	testAudience   = "cardhub-api"
	testAccessTTL  = 15 * time.Minute
	testRefreshTTL = 30 * 24 * time.Hour
)

// ============================================================
// Fake User Repository
// ============================================================

type usersFake struct {
	user             User
	createErr        error
	getErr           error
	insertedHash     string
	insertedEmail    string
	insertedUsername string
}

func (f *usersFake) Create(
	_ context.Context,
	email, username, hash string,
) (User, error) {
	f.insertedHash = hash
	f.insertedEmail = email
	f.insertedUsername = username

	if f.createErr != nil {
		return User{}, f.createErr
	}

	u := User{
		ID:           userID,
		Email:        email,
		Username:     username,
		PasswordHash: hash,
		IsActive:     true,
	}

	f.user = u
	return u, nil
}

func (f *usersFake) GetByID(
	_ context.Context,
	_ string,
) (User, error) {
	if f.getErr != nil {
		return User{}, f.getErr
	}

	return f.user, nil
}

func (f *usersFake) GetByEmail(
	_ context.Context,
	_ string,
) (User, error) {
	if f.getErr != nil {
		return User{}, f.getErr
	}

	return f.user, nil
}

// ============================================================
// Fake Session Repository
// ============================================================

type sessionsFake struct {
	mu sync.Mutex

	userID     string
	activeHash string
	expiry     time.Time
	revoked    bool

	createCalls int
	rotateCalls int
}

func (f *sessionsFake) CreateSession(
	_ context.Context,
	id, hash string,
	exp time.Time,
) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.userID = id
	f.activeHash = hash
	f.expiry = exp
	f.revoked = false
	f.createCalls++

	return nil
}

func (f *sessionsFake) RotateSession(
	_ context.Context,
	oldHash, newHash string,
	exp time.Time,
) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.rotateCalls++

	if f.revoked ||
		oldHash != f.activeHash ||
		!f.expiry.After(time.Now()) {
		return "", ErrInvalidSession
	}

	f.activeHash = newHash
	f.expiry = exp

	return f.userID, nil
}

func (f *sessionsFake) RevokeSession(
	_ context.Context,
	hash string,
) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if hash == f.activeHash {
		f.revoked = true
	}

	return nil
}

// ============================================================
// Test Setup
// ============================================================

func newTestService(
	t *testing.T,
) (*Service, *usersFake, *sessionsFake) {
	t.Helper()

	// Generate an isolated Ed25519 signing key.
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	der, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	keyPath := filepath.Join(
		t.TempDir(),
		"private.pem",
	)

	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: der,
	})

	if err := os.WriteFile(keyPath, pemData, 0600); err != nil {
		t.Fatal(err)
	}

	// Configure JWT signing independently of environment variables.
	signer, err := token.NewSigner(
		keyPath,
		token.Config{
			Issuer:    testIssuer,
			Audience:  testAudience,
			AccessTTL: testAccessTTL,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	users := &usersFake{}
	sessions := &sessionsFake{}

	svc := New(
		users,
		sessions,
		signer,
		testRefreshTTL,
	)

	return svc, users, sessions
}

// ============================================================
// Registration Tests
// ============================================================

func TestRegisterHashesNormalizesAndRedacts(t *testing.T) {
	svc, users, _ := newTestService(t)

	got, err := svc.Register(
		context.Background(),
		"  TEST@Example.com  ",
		"test_user",
		goodPassword,
	)
	if err != nil {
		t.Fatal(err)
	}

	if got.ID != userID || got.PasswordHash != "" {
		t.Fatalf("unexpected public user: %+v", got)
	}

	if users.insertedEmail != "test@example.com" ||
		users.insertedUsername != "test_user" {
		t.Fatalf(
			"unexpected normalized input: %q %q",
			users.insertedEmail,
			users.insertedUsername,
		)
	}

	if users.insertedHash == goodPassword ||
		!strings.HasPrefix(users.insertedHash, "$argon2id$") {
		t.Fatal("password not hashed with Argon2id")
	}

	valid, err := password.Verify(
		goodPassword,
		users.insertedHash,
	)
	if err != nil || !valid {
		t.Fatalf("stored hash not verifiable: %v", err)
	}
}

func TestRegisterValidationAndDuplicate(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		username string
		password string
		want     error
	}{
		{
			name:     "invalid email",
			email:    "bad-email",
			username: "user_name",
			password: goodPassword,
			want:     ErrInvalidEmail,
		},
		{
			name:     "invalid username",
			email:    "user@example.com",
			username: "bad name",
			password: goodPassword,
			want:     ErrInvalidUsername,
		},
		{
			name:     "invalid password",
			email:    "user@example.com",
			username: "user_name",
			password: "weak",
			want:     ErrInvalidPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _, _ := newTestService(t)

			_, err := svc.Register(
				context.Background(),
				tt.email,
				tt.username,
				tt.password,
			)

			if !errors.Is(err, tt.want) {
				t.Fatalf(
					"got %v, want %v",
					err,
					tt.want,
				)
			}
		})
	}

	t.Run("duplicate user", func(t *testing.T) {
		svc, users, _ := newTestService(t)

		users.createErr = ErrUserExists

		_, err := svc.Register(
			context.Background(),
			"user@example.com",
			"user_name",
			goodPassword,
		)

		if !errors.Is(err, ErrUserExists) {
			t.Fatalf("expected duplicate user error, got %v", err)
		}
	})
}

// ============================================================
// Login, Refresh and Logout Tests
// ============================================================

func TestLoginRefreshLogout(t *testing.T) {
	svc, users, sessions := newTestService(t)
	ctx := context.Background()

	_, err := svc.Register(
		ctx,
		"user@example.com",
		"user_name",
		goodPassword,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Login should create a refresh session.
	login, err := svc.Login(
		ctx,
		"USER@EXAMPLE.COM",
		goodPassword,
	)
	if err != nil {
		t.Fatal(err)
	}

	if login.AccessToken == "" ||
		login.RefreshToken == "" {
		t.Fatal("missing authentication tokens")
	}

	if sessions.userID != userID ||
		sessions.activeHash != hashRefreshToken(login.RefreshToken) ||
		sessions.createCalls != 1 {
		t.Fatal("session not stored as token hash")
	}

	if sessions.activeHash == login.RefreshToken {
		t.Fatal("stored plaintext refresh token")
	}

	// Refresh should replace the previous token.
	rotated, err := svc.RefreshToken(
		ctx,
		login.RefreshToken,
	)
	if err != nil {
		t.Fatal(err)
	}

	if rotated.AccessToken == "" ||
		rotated.RefreshToken == "" ||
		rotated.RefreshToken == login.RefreshToken {
		t.Fatal("refresh-token rotation failed")
	}

	// Old refresh token must no longer work.
	_, err = svc.RefreshToken(ctx, login.RefreshToken)
	if !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("old refresh token reusable: %v", err)
	}

	// Logout invalidates the current refresh token.
	if err := svc.Logout(ctx, rotated.RefreshToken); err != nil {
		t.Fatal(err)
	}

	// Logout is idempotent.
	if err := svc.Logout(ctx, rotated.RefreshToken); err != nil {
		t.Fatalf("logout not idempotent: %v", err)
	}

	_, err = svc.RefreshToken(ctx, rotated.RefreshToken)
	if !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("revoked token accepted: %v", err)
	}

	// Disabled accounts cannot log in.
	users.user.IsActive = false

	_, err = svc.Login(
		ctx,
		"user@example.com",
		goodPassword,
	)
	if !errors.Is(err, ErrAccountDisabled) {
		t.Fatalf("inactive login: %v", err)
	}
}

func TestLoginUsesConfiguredRefreshTTL(t *testing.T) {
	svc, _, sessions := newTestService(t)
	ctx := context.Background()

	_, err := svc.Register(
		ctx,
		"user@example.com",
		"test_user",
		goodPassword,
	)
	if err != nil {
		t.Fatal(err)
	}

	before := time.Now()

	_, err = svc.Login(
		ctx,
		"user@example.com",
		goodPassword,
	)
	if err != nil {
		t.Fatal(err)
	}

	after := time.Now()

	minExpiry := before.Add(testRefreshTTL)
	maxExpiry := after.Add(testRefreshTTL)

	if sessions.expiry.Before(minExpiry) ||
		sessions.expiry.After(maxExpiry) {
		t.Fatalf(
			"unexpected refresh expiry: got %v, expected between %v and %v",
			sessions.expiry,
			minExpiry,
			maxExpiry,
		)
	}
}

// ============================================================
// Login Rejection and GetMe Tests
// ============================================================

func TestLoginRejectionAndGetMe(t *testing.T) {
	svc, users, sessions := newTestService(t)
	ctx := context.Background()

	_, err := svc.Register(
		ctx,
		"user@example.com",
		"user_name",
		goodPassword,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Incorrect password.
	_, err = svc.Login(
		ctx,
		"user@example.com",
		"wrong-password",
	)
	if !errors.Is(err, ErrInvalidCredentials) ||
		sessions.createCalls != 0 {
		t.Fatalf("bad login: %v", err)
	}

	// Unknown user.
	users.getErr = ErrNotFound

	_, err = svc.Login(
		ctx,
		"notfound@example.com",
		goodPassword,
	)
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("unknown user: %v", err)
	}

	users.getErr = nil

	// GetMe returns the user without password hash.
	got, err := svc.GetMe(ctx, userID)
	if err != nil ||
		got.PasswordHash != "" ||
		got.ID != userID {
		t.Fatalf("GetMe: %+v %v", got, err)
	}

	// Invalid user ID.
	_, err = svc.GetMe(ctx, "invalid")
	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("invalid ID: %v", err)
	}

	// Disabled account.
	users.user.IsActive = false

	_, err = svc.GetMe(ctx, userID)
	if !errors.Is(err, ErrAccountDisabled) {
		t.Fatalf("disabled GetMe: %v", err)
	}
}

// ============================================================
// Refresh Token Validation Tests
// ============================================================

func TestRefreshTokenInputValidation(t *testing.T) {
	svc, _, sessions := newTestService(t)

	invalidTokens := []string{
		"",
		"not-a-refresh-token",
		strings.Repeat("a", 42),
	}

	for _, refreshToken := range invalidTokens {
		_, err := svc.RefreshToken(
			context.Background(),
			refreshToken,
		)

		if !errors.Is(err, ErrInvalidSession) {
			t.Fatalf(
				"token %q: expected invalid session, got %v",
				refreshToken,
				err,
			)
		}
	}

	if sessions.rotateCalls != 0 {
		t.Fatal("malformed tokens reached repository")
	}
}
