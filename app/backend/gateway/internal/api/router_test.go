package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	examplev1 "cardhub/backend/protocol/gen/example/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeClient struct{ notFound bool }

func (c fakeClient) CreateExample(_ context.Context, in *examplev1.CreateExampleRequest, _ ...grpc.CallOption) (*examplev1.CreateExampleResponse, error) {
	return &examplev1.CreateExampleResponse{Id: "123e4567-e89b-12d3-a456-426614174000", Name: in.Name}, nil
}
func (c fakeClient) GetExample(_ context.Context, in *examplev1.GetExampleRequest, _ ...grpc.CallOption) (*examplev1.GetExampleResponse, error) {
	if c.notFound {
		return nil, status.Error(codes.NotFound, "missing")
	}
	return &examplev1.GetExampleResponse{Id: in.Id, Name: "Pikachu"}, nil
}
func TestGatewayEndpointsAndOpenAPI(t *testing.T) {
	router := Router(fakeClient{}, zap.NewNop())
	for _, path := range []string{"/openapi.json", "/healthz", "/docs"} {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("GET %s: status=%d body=%s", path, rec.Code, rec.Body.String())
		}
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/openapi.json", nil))
	var spec map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &spec); err != nil {
		t.Fatal(err)
	}

	paths, ok := spec["paths"].(map[string]any)
	if !ok || paths["/v1/examples"] == nil || paths["/v1/examples/{id}"] == nil {
		t.Fatalf("missing paths: %+v", paths)
	}

	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/examples", strings.NewReader(`{"name":"Pikachu"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated || !strings.Contains(rec.Body.String(), "Pikachu") {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}
func TestGatewayMapsNotFound(t *testing.T) {
	router := Router(fakeClient{notFound: true}, zap.NewNop())
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/examples/123e4567-e89b-12d3-a456-426614174000", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}
