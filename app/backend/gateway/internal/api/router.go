package api

import (
	"net/http"
	"time"

	"cardhub/backend/gateway/internal/auth"

	examplev1 "cardhub/backend/protocol/gen/example/v1"
	identityv1 "cardhub/backend/protocol/gen/identity/v1"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func Router(
	exampleClient examplev1.ExampleServiceClient,
	identityClient identityv1.IdentityServiceClient,
	verifier *auth.Verifier,
	logger *zap.Logger,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			req *http.Request,
		) {
			started := time.Now()
			ww := middleware.NewWrapResponseWriter(
				w, req.ProtoMajor,
			)

			next.ServeHTTP(ww, req)

			logger.Info("http request",
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.Int("status", ww.Status()),
				zap.Duration("duration", time.Since(started)),
			)
		})
	})

	// Protect all routes except the explicit public allowlist.
	r.Use(auth.Middleware(verifier))

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	cfg := huma.DefaultConfig("CardHub API", "0.1.0")
	cfg.Info.Description = "CardHub authenticated REST API"

	cfg.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"BearerAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
			Description:  "Enter your access token obtained from /v1/auth/login",
		},
	}

	h := humachi.New(r, cfg)

	registerExample(h, exampleClient, logger)
	registerIdentity(h, identityClient, logger)

	return r
}
