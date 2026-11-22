package api

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"cardhub/backend/gateway/internal/auth"
	examplev1 "cardhub/backend/protocol/gen/example/v1"
	identityv1 "cardhub/backend/protocol/gen/identity/v1"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	testUserID   = "123e4567-e89b-12d3-a456-426614174000"
	testIssuer   = "cardhub-identity"
	testAudience = "cardhub-api"
)

// ============================================================
// Fake gRPC Clients
// ============================================================

type fakeClient struct {
	notFound bool
}

func (c fakeClient) CreateExample(
	_ context.Context,
	in *examplev1.CreateExampleRequest,
	_ ...grpc.CallOption,
) (*examplev1.CreateExampleResponse, error) {
	return &examplev1.CreateExampleResponse{
		Id:   testUserID,
		Name: in.Name,
	}, nil
}

func (c fakeClient) GetExample(
	_ context.Context,
	in *examplev1.GetExampleRequest,
	_ ...grpc.CallOption,
) (*examplev1.GetExampleResponse, error) {
	if c.notFound {
		return nil, status.Error(codes.NotFound, "missing")
	}

	return &examplev1.GetExampleResponse{
		Id:   in.Id,
		Name: "Pikachu",
	}, nil
}

// Embed the generated interface so tests only need
// to implement the methods they exercise.
type fakeIdentityClient struct {
	identityv1.IdentityServiceClient
}

func (c fakeIdentityClient) GetMe(
	_ context.Context,
	_ *identityv1.GetMeRequest,
	_ ...grpc.CallOption,
) (*identityv1.GetMeResponse, error) {
	return &identityv1.GetMeResponse{
		UserId:   testUserID,
		Email:    "test@example.com",
		Username: "testuser",
	}, nil
}

// ============================================================
// Test Helpers
// ============================================================

// testAuth generates a temporary Ed25519 key pair and
// a valid JWT using explicit test configuration.
func testAuth(t *testing.T) (*auth.Verifier, string) {
	t.Helper()

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	der, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		t.Fatal(err)
	}

	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	})

	keyPath := filepath.Join(
		t.TempDir(),
		"jwt-public.pem",
	)

	if err := os.WriteFile(keyPath, pemData, 0600); err != nil {
		t.Fatal(err)
	}

	// Configure JWT verification explicitly.
	verifier, err := auth.NewVerifier(
		keyPath,
		testIssuer,
		testAudience,
	)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now().UTC()

	claims := jwt.RegisteredClaims{
		Subject:   testUserID,
		Issuer:    testIssuer,
		Audience:  jwt.ClaimStrings{testAudience},
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(
			now.Add(15 * time.Minute),
		),
	}

	signedToken, err := jwt.NewWithClaims(
		jwt.SigningMethodEdDSA,
		claims,
	).SignedString(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	return verifier, signedToken
}

func testRouter(
	t *testing.T,
	example fakeClient,
) (http.Handler, string) {
	t.Helper()

	verifier, accessToken := testAuth(t)

	router := Router(
		example,
		fakeIdentityClient{},
		verifier,
		zap.NewNop(),
	)

	return router, accessToken
}

func authenticatedRequest(
	method, path, token string,
	body *strings.Reader,
) *http.Request {
	var req *http.Request

	if body == nil {
		req = httptest.NewRequest(method, path, nil)
	} else {
		req = httptest.NewRequest(method, path, body)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

// ============================================================
// Gateway and OpenAPI Tests
// ============================================================

func TestGatewayEndpointsAndOpenAPI(t *testing.T) {
	router, token := testRouter(t, fakeClient{})

	for _, path := range []string{
		"/openapi.json",
		"/healthz",
		"/docs",
	} {
		rec := httptest.NewRecorder()

		req := authenticatedRequest(
			http.MethodGet,
			path,
			token,
			nil,
		)

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf(
				"GET %s: status=%d body=%s",
				path,
				rec.Code,
				rec.Body.String(),
			)
		}
	}

	rec := httptest.NewRecorder()

	router.ServeHTTP(
		rec,
		authenticatedRequest(
			http.MethodGet,
			"/openapi.json",
			token,
			nil,
		),
	)

	var spec map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &spec); err != nil {
		t.Fatal(err)
	}

	paths, ok := spec["paths"].(map[string]any)
	if !ok ||
		paths["/v1/examples"] == nil ||
		paths["/v1/examples/{id}"] == nil ||
		paths["/v1/users/me"] == nil {
		t.Fatalf("missing paths: %+v", paths)
	}

	rec = httptest.NewRecorder()

	req := authenticatedRequest(
		http.MethodPost,
		"/v1/examples",
		token,
		strings.NewReader(`{"name":"Pikachu"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated ||
		!strings.Contains(rec.Body.String(), "Pikachu") {
		t.Fatalf(
			"status=%d body=%s",
			rec.Code,
			rec.Body.String(),
		)
	}
}

// ============================================================
// Error Handling Tests
// ============================================================

func TestGatewayMapsNotFound(t *testing.T) {
	router, token := testRouter(t, fakeClient{
		notFound: true,
	})

	rec := httptest.NewRecorder()

	req := authenticatedRequest(
		http.MethodGet,
		"/v1/examples/"+testUserID,
		token,
		nil,
	)

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"status=%d body=%s",
			rec.Code,
			rec.Body.String(),
		)
	}
}

func TestGatewayRejectsMissingToken(t *testing.T) {
	router, _ := testRouter(t, fakeClient{})

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodGet,
		"/v1/users/me",
		nil,
	)

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected 401, got %d: %s",
			rec.Code,
			rec.Body.String(),
		)
	}
}

func TestGatewayRejectsInvalidToken(t *testing.T) {
	router, _ := testRouter(t, fakeClient{})

	rec := httptest.NewRecorder()

	req := authenticatedRequest(
		http.MethodGet,
		"/v1/users/me",
		"invalid-token",
		nil,
	)

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected 401, got %d: %s",
			rec.Code,
			rec.Body.String(),
		)
	}
}

// ============================================================
// Authenticated User Tests
// ============================================================

func TestGatewayGetMe(t *testing.T) {
	router, token := testRouter(t, fakeClient{})

	rec := httptest.NewRecorder()

	req := authenticatedRequest(
		http.MethodGet,
		"/v1/users/me",
		token,
		nil,
	)

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200, got %d: %s",
			rec.Code,
			rec.Body.String(),
		)
	}

	var user struct {
		UserID   string `json:"user_id"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &user); err != nil {
		t.Fatal(err)
	}

	if user.UserID != testUserID {
		t.Fatalf(
			"unexpected user ID: %s",
			user.UserID,
		)
	}
}
