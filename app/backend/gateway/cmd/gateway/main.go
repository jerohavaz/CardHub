package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cardhub/backend/foundation/config"
	"cardhub/backend/foundation/logging"
	"cardhub/backend/gateway/internal/api"
	"cardhub/backend/gateway/internal/auth"
	"cardhub/backend/gateway/internal/client"

	"go.uber.org/zap"
)

func main() {
	logger, err := logging.New(config.Env("APP_ENV", "development") == "development")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	rpc, conn, err := client.DialExample(config.Env("EXAMPLE_GRPC_ADDR", "localhost:9000"))
	if err != nil {
		logger.Fatal("dial example gRPC", zap.Error(err))
	}
	defer conn.Close()

	identityRPC, identityConn, err := client.DialIdentity(
		config.Env("IDENTITY_GRPC_ADDR", "localhost:9001"),
	)
	if err != nil {
		logger.Fatal("dial identity gRPC", zap.Error(err))
	}
	defer identityConn.Close()

	verifier, err := auth.NewVerifier(
		config.Env("JWT_PUBLIC_KEY_PATH", "/run/keys/jwt-public.pem"),
	)
	if err != nil {
		logger.Fatal("load JWT public key", zap.Error(err))
	}

	server := &http.Server{
		Addr: config.Env("HTTP_ADDR", ":8080"),
		Handler: api.Router(
			rpc,
			identityRPC,
			verifier,
			logger,
		),
		ReadHeaderTimeout: 5 * time.Second,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	done := make(chan error, 1)
	go func() { done <- server.ListenAndServe() }()

	logger.Info("http listening", zap.String("address", server.Addr))

	select {
	case err := <-done:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("serve http", zap.Error(err))
		}
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutCtx); err != nil {
			logger.Error("shutdown http", zap.Error(err))
		}
	}
}
