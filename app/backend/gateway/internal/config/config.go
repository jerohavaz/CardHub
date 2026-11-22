package config

import (
	"errors"
	"strings"

	"cardhub/backend/foundation/config"
)

type Config struct {
	AppEnv string

	HTTPAddr         string
	ExampleGRPCAddr  string
	IdentityGRPCAddr string

	JWTPublicKeyPath string
	JWTIssuer        string
	JWTAudience      string
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv: config.Env("APP_ENV", "development"),

		HTTPAddr:         config.Env("HTTP_ADDR", ":8080"),
		ExampleGRPCAddr:  config.Env("EXAMPLE_GRPC_ADDR", "localhost:9000"),
		IdentityGRPCAddr: config.Env("IDENTITY_GRPC_ADDR", "localhost:9001"),

		JWTPublicKeyPath: config.Env("JWT_PUBLIC_KEY_PATH", "/run/keys/jwt-public.pem"),
		JWTIssuer:        config.Env("JWT_ISSUER", "cardhub-identity"),
		JWTAudience:      config.Env("JWT_AUDIENCE", "cardhub-api"),
	}

	if strings.TrimSpace(cfg.JWTIssuer) == "" ||
		strings.TrimSpace(cfg.JWTAudience) == "" {
		return Config{}, errors.New("JWT issuer and audience must not be empty")
	}

	return cfg, nil
}
