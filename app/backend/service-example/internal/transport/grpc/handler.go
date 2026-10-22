// Package grpc adapts generated protobuf messages to application operations.
package grpc

import (
	"context"
	"errors"

	examplev1 "cardhub/backend/protocol/gen/example/v1"
	"cardhub/backend/service-example/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	examplev1.UnimplementedExampleServiceServer
	svc *service.Service
	log *zap.Logger
}

func New(svc *service.Service, log *zap.Logger) *Handler { return &Handler{svc: svc, log: log} }
func (h *Handler) CreateExample(ctx context.Context, req *examplev1.CreateExampleRequest) (*examplev1.CreateExampleResponse, error) {
	got, err := h.svc.Create(ctx, req.GetName())
	if err != nil {
		return nil, h.mapError(err)
	}
	return &examplev1.CreateExampleResponse{Id: got.ID, Name: got.Name}, nil
}
func (h *Handler) GetExample(ctx context.Context, req *examplev1.GetExampleRequest) (*examplev1.GetExampleResponse, error) {
	got, err := h.svc.Get(ctx, req.GetId())
	if err != nil {
		return nil, h.mapError(err)
	}
	return &examplev1.GetExampleResponse{Id: got.ID, Name: got.Name}, nil
}
func (h *Handler) mapError(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidName), errors.Is(err, service.ErrInvalidID):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrNotFound):
		return status.Error(codes.NotFound, "example not found")
	default:
		h.log.Error("example operation failed", zap.Error(err))
		return status.Error(codes.Internal, "internal server error")
	}
}
