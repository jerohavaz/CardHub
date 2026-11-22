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

const userID = "123e4567-e89b-12d3-a456-426614174000"
const goodPassword = "SecurePassword123!"

type usersFake struct {
	user             User
	createErr        error
	getErr           error
	insertedHash     string
	insertedEmail    string
	insertedUsername string
}

func (f *usersFake) Create(_ context.Context, email, username, hash string) (User, error) {
	f.insertedHash, f.insertedEmail, f.insertedUsername = hash, email, username
	if f.createErr != nil {
		return User{}, f.createErr
	}
	u := User{ID: userID, Email: email, Username: username, PasswordHash: hash, IsActive: true}
	f.user = u
	return u, nil
}
func (f *usersFake) GetByID(_ context.Context, _ string) (User, error) {
	if f.getErr != nil {
		return User{}, f.getErr
	}
	return f.user, nil
}
func (f *usersFake) GetByEmail(_ context.Context, _ string) (User, error) {
	if f.getErr != nil {
		return User{}, f.getErr
	}
	return f.user, nil
}

type sessionsFake struct {
	mu          sync.Mutex
	userID      string
	activeHash  string
	expiry      time.Time
	revoked     bool
	createCalls int
	rotateCalls int
}

func (f *sessionsFake) CreateSession(_ context.Context, id, hash string, exp time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.userID = id
	f.activeHash = hash
	f.expiry = exp
	f.revoked = false
	f.createCalls++
	return nil
}
func (f *sessionsFake) RotateSession(_ context.Context, oldHash, newHash string, exp time.Time) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.rotateCalls++
	if f.revoked || oldHash != f.activeHash || !f.expiry.After(time.Now()) {
		return "", ErrInvalidSession
	}
	f.activeHash = newHash
	f.expiry = exp
	return f.userID, nil
}
func (f *sessionsFake) RevokeSession(_ context.Context, hash string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if hash == f.activeHash {
		f.revoked = true
	}
	return nil
}
func newTestService(t *testing.T) (*Service, *usersFake, *sessionsFake) {
	t.Helper()
	_, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "private.pem")
	if err := os.WriteFile(path, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), 0600); err != nil {
		t.Fatal(err)
	}
	signer, err := token.NewSigner(path)
	if err != nil {
		t.Fatal(err)
	}
	u := &usersFake{}
	s := &sessionsFake{}
	return New(u, s, signer), u, s
}
func TestRegisterHashesNormalizesAndRedacts(t *testing.T) {
	svc, u, _ := newTestService(t)
	got, err := svc.Register(context.Background(), "  TEST@Example.com  ", "test_user", goodPassword)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != userID || got.PasswordHash != "" {
		t.Fatalf("unexpected public user: %+v", got)
	}
	if u.insertedEmail != "test@example.com" || u.insertedUsername != "test_user" {
		t.Fatalf("unexpected normalized input: %q %q", u.insertedEmail, u.insertedUsername)
	}
	if u.insertedHash == goodPassword || !strings.HasPrefix(u.insertedHash, "$argon2id$") {
		t.Fatal("password not hashed with Argon2id")
	}
	ok, err := password.Verify(goodPassword, u.insertedHash)
	if err != nil || !ok {
		t.Fatalf("stored hash not verifiable: %v", err)
	}
}
func TestRegisterValidationAndDuplicate(t *testing.T) {
	tests := []struct {
		name, email, username, pw string
		want                      error
	}{
		{"email", "bad-email", "user_name", goodPassword, ErrInvalidEmail},
		{"username", "user@example.com", "bad name", goodPassword, ErrInvalidUsername},
		{"password", "user@example.com", "user_name", "weak", ErrInvalidPassword},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _, _ := newTestService(t)
			_, err := svc.Register(context.Background(), tt.email, tt.username, tt.pw)
			if !errors.Is(err, tt.want) {
				t.Fatalf("got %v want %v", err, tt.want)
			}
		})
	}
	svc, u, _ := newTestService(t)
	u.createErr = ErrUserExists
	_, err := svc.Register(context.Background(), "user@example.com", "user_name", goodPassword)
	if !errors.Is(err, ErrUserExists) {
		t.Fatalf("duplicate: %v", err)
	}
}
func TestLoginRefreshLogout(t *testing.T) {
	svc, u, sessions := newTestService(t)
	_, err := svc.Register(context.Background(), "user@example.com", "user_name", goodPassword)
	if err != nil {
		t.Fatal(err)
	}
	login, err := svc.Login(context.Background(), "USER@EXAMPLE.COM", goodPassword)
	if err != nil {
		t.Fatal(err)
	}
	if login.AccessToken == "" || login.RefreshToken == "" {
		t.Fatal("missing tokens")
	}
	if sessions.userID != userID || sessions.activeHash != hashRefreshToken(login.RefreshToken) || sessions.createCalls != 1 {
		t.Fatal("session not stored as token hash")
	}
	if sessions.activeHash == login.RefreshToken {
		t.Fatal("stored plaintext refresh token")
	}
	rotated, err := svc.RefreshToken(context.Background(), login.RefreshToken)
	if err != nil {
		t.Fatal(err)
	}
	if rotated.AccessToken == "" || rotated.RefreshToken == login.RefreshToken {
		t.Fatal("rotation failed")
	}
	_, err = svc.RefreshToken(context.Background(), login.RefreshToken)
	if !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("old token reusable: %v", err)
	}
	if err := svc.Logout(context.Background(), rotated.RefreshToken); err != nil {
		t.Fatal(err)
	}
	if err := svc.Logout(context.Background(), rotated.RefreshToken); err != nil {
		t.Fatalf("logout not idempotent: %v", err)
	}
	_, err = svc.RefreshToken(context.Background(), rotated.RefreshToken)
	if !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("revoked token accepted: %v", err)
	}
	u.user.IsActive = false
	_, err = svc.Login(context.Background(), "user@example.com", goodPassword)
	if !errors.Is(err, ErrAccountDisabled) {
		t.Fatalf("inactive login: %v", err)
	}
}
func TestLoginRejectionAndGetMe(t *testing.T) {
	svc, u, sessions := newTestService(t)
	_, err := svc.Register(context.Background(), "user@example.com", "user_name", goodPassword)
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.Login(context.Background(), "user@example.com", "wrong-password")
	if !errors.Is(err, ErrInvalidCredentials) || sessions.createCalls != 0 {
		t.Fatalf("bad login: %v", err)
	}
	u.getErr = ErrNotFound
	_, err = svc.Login(context.Background(), "notfound@example.com", goodPassword)
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("unknown user: %v", err)
	}
	u.getErr = nil
	got, err := svc.GetMe(context.Background(), userID)
	if err != nil || got.PasswordHash != "" || got.ID != userID {
		t.Fatalf("GetMe: %+v %v", got, err)
	}
	_, err = svc.GetMe(context.Background(), "invalid")
	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("invalid ID: %v", err)
	}
	u.user.IsActive = false
	_, err = svc.GetMe(context.Background(), userID)
	if !errors.Is(err, ErrAccountDisabled) {
		t.Fatalf("disabled GetMe: %v", err)
	}
}
func TestRefreshTokenInputValidation(t *testing.T) {
	svc, _, sessions := newTestService(t)
	for _, s := range []string{"", "not-a-refresh-token", strings.Repeat("a", 42)} {
		_, err := svc.RefreshToken(context.Background(), s)
		if !errors.Is(err, ErrInvalidSession) {
			t.Fatalf("token %q: %v", s, err)
		}
	}
	if sessions.rotateCalls != 0 {
		t.Fatal("malformed tokens reached repository")
	}
}
