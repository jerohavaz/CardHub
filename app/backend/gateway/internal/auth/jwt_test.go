package auth

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTVerifier(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "public.pem")
	if err := os.WriteFile(path, pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}), 0600); err != nil {
		t.Fatal(err)
	}
	verifier, err := NewVerifier(path)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	cases := []struct {
		name   string
		claims jwt.RegisteredClaims
		method jwt.SigningMethod
		wantOK bool
	}{
		{"valid", jwt.RegisteredClaims{Subject: "user-123", Issuer: issuer, Audience: jwt.ClaimStrings{audience}, IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour))}, jwt.SigningMethodEdDSA, true},
		{"expired", jwt.RegisteredClaims{Subject: "user-123", Issuer: issuer, Audience: jwt.ClaimStrings{audience}, IssuedAt: jwt.NewNumericDate(now.Add(-2 * time.Hour)), ExpiresAt: jwt.NewNumericDate(now.Add(-time.Hour))}, jwt.SigningMethodEdDSA, false},
		{"wrong issuer", jwt.RegisteredClaims{Subject: "user-123", Issuer: "other", Audience: jwt.ClaimStrings{audience}, ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour))}, jwt.SigningMethodEdDSA, false},
		{"wrong audience", jwt.RegisteredClaims{Subject: "user-123", Issuer: issuer, Audience: jwt.ClaimStrings{"other"}, ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour))}, jwt.SigningMethodEdDSA, false},
		{"missing subject", jwt.RegisteredClaims{Issuer: issuer, Audience: jwt.ClaimStrings{audience}, ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour))}, jwt.SigningMethodEdDSA, false},
		{"missing expiration", jwt.RegisteredClaims{Subject: "user-123", Issuer: issuer, Audience: jwt.ClaimStrings{audience}}, jwt.SigningMethodEdDSA, false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			signed, err := jwt.NewWithClaims(tt.method, tt.claims).SignedString(priv)
			if err != nil {
				t.Fatal(err)
			}
			sub, err := verifier.Verify(signed)
			if (err == nil) != tt.wantOK {
				t.Fatalf("expected valid=%t got sub=%q err=%v", tt.wantOK, sub, err)
			}
			if tt.wantOK && sub != "user-123" {
				t.Fatalf("subject=%q", sub)
			}
		})
	}
	otherPub, otherPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	_ = otherPub
	signed, err := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.RegisteredClaims{
		Subject: "user-123", Issuer: issuer, Audience: jwt.ClaimStrings{audience}, ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
	}).SignedString(otherPriv)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = verifier.Verify(signed); err == nil {
		t.Fatal("accepted wrong signing key")
	}
	if _, err = verifier.Verify("not-a-jwt"); err == nil {
		t.Fatal("accepted malformed JWT")
	}
}
