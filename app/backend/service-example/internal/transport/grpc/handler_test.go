package grpc

import (
	"context"
	"testing"

	examplev1 "cardhub/backend/protocol/gen/example/v1"
	"cardhub/backend/service-example/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type stubRepo struct{}

func (stubRepo) Create(_ context.Context, n string) (service.Example, error) {
	return service.Example{ID: "id", Name: n}, nil
}
func (stubRepo) Get(_ context.Context, _ string) (service.Example, error) {
	return service.Example{}, service.ErrNotFound
}
func TestGrpcErrorMapping(t *testing.T) {
	h := New(service.New(stubRepo{}), zap.NewNop())
	_, err := h.CreateExample(context.Background(), &examplev1.CreateExampleRequest{Name: " "})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("got %v", err)
	}
	_, err = h.GetExample(context.Background(), &examplev1.GetExampleRequest{Id: "123e4567-e89b-12d3-a456-426614174000"})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("got %v", err)
	}
}
