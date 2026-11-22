package api

import (
	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func publicError(err error, logger *zap.Logger) error {
	switch status.Code(err) {
	case codes.InvalidArgument:
		return huma.Error400BadRequest(
			status.Convert(err).Message(),
		)
	case codes.AlreadyExists:
		return huma.Error409Conflict("resource already exists")
	case codes.Unauthenticated:
		return huma.Error401Unauthorized("unauthorized")
	case codes.PermissionDenied:
		return huma.Error403Forbidden("forbidden")
	case codes.NotFound:
		return huma.Error404NotFound("resource not found")
	case codes.DeadlineExceeded:
		return huma.Error504GatewayTimeout("service timed out")
	case codes.Unavailable:
		return huma.Error503ServiceUnavailable("service unavailable")
	default:
		logger.Error("backend gRPC error", zap.Error(err))
		return huma.Error502BadGateway("backend request failed")
	}
}
