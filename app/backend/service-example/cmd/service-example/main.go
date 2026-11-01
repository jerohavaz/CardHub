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
	examplev1 "cardhub/backend/protocol/gen/example/v1"
	repository "cardhub/backend/service-example/internal/postgres"
	"cardhub/backend/service-example/internal/service"
	transport "cardhub/backend/service-example/internal/transport/grpc"
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

	pool, err := foundationpostgres.Open(ctx, url)
	if err != nil {
		logger.Fatal("postgres", zap.Error(err))
	}
	defer pool.Close()

	svc := service.New(repository.New(pool))
	server := grpcserver.New(logger)
	examplev1.RegisterExampleServiceServer(server, transport.New(svc, logger))
	if err := grpcserver.Run(ctx, config.Env("GRPC_ADDR", ":9000"), server, logger); err != nil {
		logger.Fatal("grpc server", zap.Error(err))
	}
}
