package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cardhub/backend/foundation/grpcserver"
	"cardhub/backend/foundation/logging"
	foundationpostgres "cardhub/backend/foundation/postgres"
	examplev1 "cardhub/backend/protocol/gen/example/v1"
	exampleconfig "cardhub/backend/service-example/internal/config"
	repository "cardhub/backend/service-example/internal/postgres"
	"cardhub/backend/service-example/internal/service"
	transport "cardhub/backend/service-example/internal/transport/grpc"

	"go.uber.org/zap"
)

func main() {
	// Load and validate configuration.
	cfg, err := exampleconfig.Load()
	if err != nil {
		log.Fatal("configuration: ", err)
	}

	// Initialize logging.
	logger, err := logging.New(cfg.AppEnv == "development")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	// Handle shutdown signals.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// Connect to PostgreSQL.
	pool, err := foundationpostgres.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("postgres", zap.Error(err))
	}
	defer pool.Close()

	// Initialize repository and business logic.
	repo := repository.New(pool)
	svc := service.New(repo)

	// Register gRPC handlers.
	server := grpcserver.New(logger)

	examplev1.RegisterExampleServiceServer(
		server,
		transport.New(svc, logger),
	)

	// Start the gRPC server.
	if err := grpcserver.Run(
		ctx,
		cfg.GRPCAddr,
		server,
		logger,
	); err != nil {
		logger.Fatal("grpc server", zap.Error(err))
	}
}
