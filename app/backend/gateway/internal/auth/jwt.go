package auth

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

const (
	issuer   = "cardhub-identity"
	audience = "cardhub-api"
)

type Verifier struct {
	publicKey ed25519.PublicKey
}

func NewVerifier(path string) (*Verifier, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read JWT public key: %w", err)
	}

	block, rest := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" ||
		len(rest) != 0 {
		return nil, errors.New("invalid public key PEM")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	publicKey, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("expected Ed25519 public key")
	}

	return &Verifier{publicKey: publicKey}, nil
}

func (v *Verifier) Verify(raw string) (string, error) {
	claims := &jwt.RegisteredClaims{}

	parsed, err := jwt.ParseWithClaims(
		raw,
		claims,
		func(_ *jwt.Token) (any, error) {
			return v.publicKey, nil
		},
		jwt.WithValidMethods([]string{"EdDSA"}),
		jwt.WithIssuer(issuer),
		jwt.WithAudience(audience),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithLeeway(5_000_000_000), // 5 seconds
	)

	if err != nil || !parsed.Valid {
		return "", errors.New("invalid access token")
	}

	if claims.Subject == "" {
		return "", errors.New("missing JWT subject")
	}

	return claims.Subject, nil
}
