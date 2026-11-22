package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	identityv1 "cardhub/backend/protocol/gen/identity/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type capturingIdentityClient struct {
	identityv1.IdentityServiceClient
	forwardedID string
}

func (c *capturingIdentityClient) GetMe(ctx context.Context, _ *identityv1.GetMeRequest, _ ...grpc.CallOption) (*identityv1.GetMeResponse, error) {
	md, _ := metadata.FromOutgoingContext(ctx)
	ids := md.Get("x-user-id")
	if len(ids) == 1 {
		c.forwardedID = ids[0]
	}
	return &identityv1.GetMeResponse{UserId: c.forwardedID, Email: "test@example.com", Username: "testuser"}, nil
}
func TestGetMeForwardsVerifiedIdentityNotClientHeader(t *testing.T) {
	verifier, jwtToken := testAuth(t)
	fake := &capturingIdentityClient{}
	router := Router(fakeClient{}, fake, verifier, zap.NewNop())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("X-User-Id", "attacker-chosen-user")
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("got %d: %s", rec.Code, rec.Body.String())
	}
	if fake.forwardedID != testUserID {
		t.Fatalf("forged identity forwarded: %q", fake.forwardedID)
	}
}
func TestPublicDocumentationDoesNotRequireBearer(t *testing.T) {
	verifier, _ := testAuth(t)
	router := Router(fakeClient{}, fakeIdentityClient{}, verifier, zap.NewNop())
	for _, path := range []string{"/healthz", "/openapi.json", "/docs"} {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusOK {
			t.Errorf("GET %s: status=%d body=%s", path, rec.Code, rec.Body.String())
		}
	}
}
func TestJWTMiddlewareRejectsSpoofedIdentityWithoutToken(t *testing.T) {
	verifier, _ := testAuth(t)
	router := Router(fakeClient{}, fakeIdentityClient{}, verifier, zap.NewNop())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/users/me", nil)
	req.Header.Set("X-User-Id", "attacker-chosen-user")
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("got %d", rec.Code)
	}
}
