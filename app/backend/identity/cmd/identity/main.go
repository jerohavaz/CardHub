package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cardhub/backend/foundation/config"
	"cardhub/backend/foundation/grpcserver"
	"cardhub/backend/foundation/logging"
	foundationpostgres "cardhub/backend/foundation/postgres"
	repository "cardhub/backend/identity/internal/postgres"
	"cardhub/backend/identity/internal/service"
	"cardhub/backend/identity/internal/token"
	transport "cardhub/backend/identity/internal/transport/grpc"
	identityv1 "cardhub/backend/protocol/gen/identity/v1"

	"go.uber.org/zap"
)

func main() {
	logger, err := logging.New(config.Env("APP_ENV", "development") == "development")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync() // stdout sync is a no-op or unsupported on some platforms.

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	url, err := config.Required("DATABASE_URL")
	if err != nil {
		logger.Fatal("configuration", zap.Error(err))
	}

	keyPath, err := config.Required("JWT_KEY_PATH")
	if err != nil {
		logger.Fatal("configuration", zap.Error(err))
	}

	pool, err := foundationpostgres.Open(ctx, url)
	if err != nil {
		logger.Fatal("postgres", zap.Error(err))
	}
	defer pool.Close()

	signer, err := token.NewSigner(keyPath)
	if err != nil {
		logger.Fatal("signer", zap.Error(err))
	}

	svc := service.New(repository.NewUserRepository(pool), repository.NewSessionRepository(pool), signer)
	server := grpcserver.New(logger)
	identityv1.RegisterIdentityServiceServer(server, transport.New(svc, logger))
	if err := grpcserver.Run(ctx, config.Env("GRPC_ADDR", ":9001"), server, logger); err != nil {
		logger.Fatal("grpc server", zap.Error(err))
	}
}
