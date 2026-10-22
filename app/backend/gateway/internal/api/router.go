// Package api owns the public HTTP contract and OpenAPI generation.
package api

import (
	"net/http"
	"time"

	examplev1 "cardhub/backend/protocol/gen/example/v1"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func Router(client examplev1.ExampleServiceClient, logger *zap.Logger) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			started := time.Now()
			ww := middleware.NewWrapResponseWriter(w, req.ProtoMajor)
			next.ServeHTTP(ww, req)
			logger.Info("http request", zap.String("method", req.Method), zap.String("path", req.URL.Path), zap.Int("status", ww.Status()), zap.Duration("duration", time.Since(started)))
		})
	})

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	cfg := huma.DefaultConfig("CardHub scaffold", "0.1.0")
	cfg.Info.Description = "Minimal HTTP gateway to a gRPC example service. Development only; authentication is not implemented."
	h := humachi.New(r, cfg)
	registerExample(h, client, logger)
	return r
}
