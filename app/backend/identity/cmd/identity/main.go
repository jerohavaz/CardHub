package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	accessTTL, err := config.DurationEnv(
		"JWT_ACCESS_TTL",
		15*time.Minute,
	)
	if err != nil {
		logger.Fatal("invalid JWT access TTL", zap.Error(err))
	}

	signer, err := token.NewSigner(
		keyPath,
		token.Config{
			Issuer:    config.Env("JWT_ISSUER", "cardhub-identity"),
			Audience:  config.Env("JWT_AUDIENCE", "cardhub-api"),
			AccessTTL: accessTTL,
		},
	)

	refreshTTL, err := config.DurationEnv("JWT_REFRESH_TTL", 720*24*time.Minute)
	if err != nil {
		logger.Fatal("invalid JWT refresh TTL", zap.Error(err))
	}

	svc := service.New(repository.NewUserRepository(pool), repository.NewSessionRepository(pool), signer, refreshTTL)
	server := grpcserver.New(logger)
	identityv1.RegisterIdentityServiceServer(server, transport.New(svc, logger))
	if err := grpcserver.Run(ctx, config.Env("GRPC_ADDR", ":9001"), server, logger); err != nil {
		logger.Fatal("grpc server", zap.Error(err))
	}
}
