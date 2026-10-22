package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// New configures shared logging, panic recovery, gRPC health, and reflection.
// Plaintext transport is intended only for private local Docker networking.
func New(log *zap.Logger) *grpc.Server {
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(recoveryInterceptor(log), loggingInterceptor(log)))
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	reflection.Register(server)
	return server
}

// Run serves until cancellation and bounds graceful shutdown to ten seconds.
func Run(ctx context.Context, addr string, server *grpc.Server, log *zap.Logger) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen gRPC %s: %w", addr, err)
	}

	log.Info("grpc listening", zap.String("address", addr))

	done := make(chan error, 1)
	go func() { done <- server.Serve(lis) }()

	select {
	case err := <-done:
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return fmt.Errorf("serve grpc: %w", err)
		}
		return nil
	case <-ctx.Done():
	}

	stopped := make(chan struct{})
	go func() { server.GracefulStop(); close(stopped) }()

	select {
	case <-stopped:
		return nil
	case <-time.After(10 * time.Second):
		server.Stop()
		<-stopped
		return fmt.Errorf("grpc graceful shutdown timed out")
	}
}
