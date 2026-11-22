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

	"cardhub/backend/foundation/logging"
	"cardhub/backend/gateway/internal/api"
	"cardhub/backend/gateway/internal/auth"
	"cardhub/backend/gateway/internal/client"
	"cardhub/backend/gateway/internal/config"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logging.New(cfg.AppEnv == "development")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	// gRPC clients.
	rpc, conn, err := client.DialExample(cfg.ExampleGRPCAddr)
	if err != nil {
		logger.Fatal("dial example gRPC", zap.Error(err))
	}
	defer conn.Close()

	identityRPC, identityConn, err := client.DialIdentity(cfg.IdentityGRPCAddr)
	if err != nil {
		logger.Fatal("dial identity gRPC", zap.Error(err))
	}
	defer identityConn.Close()

	// JWT verification.
	verifier, err := auth.NewVerifier(
		cfg.JWTPublicKeyPath,
		cfg.JWTIssuer,
		cfg.JWTAudience,
	)
	if err != nil {
		logger.Fatal("load JWT public key", zap.Error(err))
	}

	// HTTP server.
	server := &http.Server{
		Addr: cfg.HTTPAddr,
		Handler: api.Router(
			rpc,
			identityRPC,
			verifier,
			logger,
		),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Graceful shutdown.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	done := make(chan error, 1)
	go func() {
		done <- server.ListenAndServe()
	}()

	logger.Info("http listening", zap.String("address", server.Addr))

	select {
	case err := <-done:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("serve http", zap.Error(err))
		}

	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(shutCtx); err != nil {
			logger.Error("shutdown http", zap.Error(err))
		}
	}
}
