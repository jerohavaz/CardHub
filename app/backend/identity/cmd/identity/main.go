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
	"cardhub/backend/identity/internal/config"
	repository "cardhub/backend/identity/internal/postgres"
	"cardhub/backend/identity/internal/service"
	"cardhub/backend/identity/internal/token"
	transport "cardhub/backend/identity/internal/transport/grpc"
	identityv1 "cardhub/backend/protocol/gen/identity/v1"

	"go.uber.org/zap"
)

func main() {
	// Load and validate configuration.
	cfg, err := config.Load()
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

	// Initialize JWT signer.
	signer, err := token.NewSigner(
		cfg.JWTKeyPath,
		token.Config{
			Issuer:    cfg.JWTIssuer,
			Audience:  cfg.JWTAudience,
			AccessTTL: cfg.JWTAccessTTL,
		},
	)
	if err != nil {
		logger.Fatal("JWT signer", zap.Error(err))
	}

	// Initialize repositories and business logic.
	users := repository.NewUserRepository(pool)
	sessions := repository.NewSessionRepository(pool)

	svc := service.New(
		users,
		sessions,
		signer,
		cfg.JWTRefreshTTL,
	)

	// Register gRPC handlers.
	server := grpcserver.New(logger)

	identityv1.RegisterIdentityServiceServer(
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
