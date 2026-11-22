package token

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	issuer    = "cardhub-identity"
	audience  = "cardhub-api"
	accessTTL = 15 * time.Minute
)

type Signer struct {
	privateKey ed25519.PrivateKey
}

func NewSigner(keyPath string) (*Signer, error) {
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("read JWT private key: %w", err)
	}

	block, rest := pem.Decode(data)
	if block == nil || block.Type != "PRIVATE KEY" || len(rest) != 0 {
		return nil, errors.New("invalid private key PEM")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	privateKey, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("expected Ed25519 private key")
	}

	return &Signer{privateKey: privateKey}, nil
}

func (s *Signer) Issue(userID string) (string, error) {
	if userID == "" {
		return "", errors.New("user ID is required")
	}

	now := time.Now().UTC()

	claims := jwt.RegisteredClaims{
		Subject:   userID,
		Issuer:    issuer,
		Audience:  jwt.ClaimStrings{audience},
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(accessTTL)),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	signed, err := t.SignedString(s.privateKey)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return signed, nil
}
