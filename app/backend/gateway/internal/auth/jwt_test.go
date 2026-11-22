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

const (
	testIssuer   = "cardhub-identity"
	testAudience = "cardhub-api"
	testUserID   = "123e4567-e89b-12d3-a456-426614174000"
)

func TestJWTVerifier(t *testing.T) {
	// Generate temporary Ed25519 keys.
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}

	keyPath := filepath.Join(t.TempDir(), "public.pem")
	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	})

	if err := os.WriteFile(keyPath, pemData, 0600); err != nil {
		t.Fatal(err)
	}

	// Pass the configured issuer and audience.
	verifier, err := NewVerifier(
		keyPath,
		testIssuer,
		testAudience,
	)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now().UTC()

	tests := []struct {
		name   string
		claims jwt.RegisteredClaims
		wantOK bool
	}{
		{
			name: "valid token",
			claims: jwt.RegisteredClaims{
				Subject:   testUserID,
				Issuer:    testIssuer,
				Audience:  jwt.ClaimStrings{testAudience},
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			},
			wantOK: true,
		},
		{
			name: "expired token",
			claims: jwt.RegisteredClaims{
				Subject:   testUserID,
				Issuer:    testIssuer,
				Audience:  jwt.ClaimStrings{testAudience},
				IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)),
				ExpiresAt: jwt.NewNumericDate(now.Add(-time.Hour)),
			},
		},
		{
			name: "wrong issuer",
			claims: jwt.RegisteredClaims{
				Subject:   testUserID,
				Issuer:    "other-issuer",
				Audience:  jwt.ClaimStrings{testAudience},
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			},
		},
		{
			name: "wrong audience",
			claims: jwt.RegisteredClaims{
				Subject:   testUserID,
				Issuer:    testIssuer,
				Audience:  jwt.ClaimStrings{"other-audience"},
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			},
		},
		{
			name: "missing subject",
			claims: jwt.RegisteredClaims{
				Issuer:    testIssuer,
				Audience:  jwt.ClaimStrings{testAudience},
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			},
		},
		{
			name: "missing expiration",
			claims: jwt.RegisteredClaims{
				Subject:  testUserID,
				Issuer:   testIssuer,
				Audience: jwt.ClaimStrings{testAudience},
			},
		},
		{
			name: "not valid yet",
			claims: jwt.RegisteredClaims{
				Subject:   testUserID,
				Issuer:    testIssuer,
				Audience:  jwt.ClaimStrings{testAudience},
				NotBefore: jwt.NewNumericDate(now.Add(time.Hour)),
				ExpiresAt: jwt.NewNumericDate(now.Add(2 * time.Hour)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signed, err := jwt.NewWithClaims(
				jwt.SigningMethodEdDSA,
				tt.claims,
			).SignedString(priv)
			if err != nil {
				t.Fatal(err)
			}

			subject, err := verifier.Verify(signed)

			if (err == nil) != tt.wantOK {
				t.Fatalf(
					"want valid=%t, got subject=%q, err=%v",
					tt.wantOK,
					subject,
					err,
				)
			}

			if tt.wantOK && subject != testUserID {
				t.Fatalf(
					"unexpected subject: %q",
					subject,
				)
			}
		})
	}
}

func TestJWTVerifierRejectsWrongKey(t *testing.T) {
	// Generate the trusted key.
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}

	keyPath := filepath.Join(t.TempDir(), "public.pem")
	if err := os.WriteFile(
		keyPath,
		pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: der,
		}),
		0600,
	); err != nil {
		t.Fatal(err)
	}

	verifier, err := NewVerifier(
		keyPath,
		testIssuer,
		testAudience,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Sign with an unrelated private key.
	_, otherPrivateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	signed, err := jwt.NewWithClaims(
		jwt.SigningMethodEdDSA,
		jwt.RegisteredClaims{
			Subject:   testUserID,
			Issuer:    testIssuer,
			Audience:  jwt.ClaimStrings{testAudience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	).SignedString(otherPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := verifier.Verify(signed); err == nil {
		t.Fatal("accepted token signed with wrong key")
	}
}

func TestJWTVerifierRejectsMalformedToken(t *testing.T) {
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}

	keyPath := filepath.Join(t.TempDir(), "public.pem")
	if err := os.WriteFile(
		keyPath,
		pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: der,
		}),
		0600,
	); err != nil {
		t.Fatal(err)
	}

	verifier, err := NewVerifier(
		keyPath,
		testIssuer,
		testAudience,
	)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := verifier.Verify("not-a-jwt"); err == nil {
		t.Fatal("accepted malformed JWT")
	}
}
