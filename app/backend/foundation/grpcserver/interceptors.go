// Package grpcserver installs shared unary middleware for internal services.
// Adapted from TimBRA: logging and panic recovery, without rate limits or mandatory TLS.
package grpcserver

import (
	"context"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func loggingInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		result, err := handler(ctx, req)
		fields := []zap.Field{zap.String("method", info.FullMethod), zap.String("code", status.Code(err).String()), zap.Duration("duration", time.Since(start))}

		if err != nil {
			log.Warn("grpc request", append(fields, zap.Error(err))...)
		} else {
			log.Info("grpc request", fields...)
		}

		return result, err
	}
}
func recoveryInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (response any, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("grpc panic", zap.String("method", info.FullMethod), zap.Any("panic", r), zap.ByteString("stack", debug.Stack()))
				err = status.Error(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}
