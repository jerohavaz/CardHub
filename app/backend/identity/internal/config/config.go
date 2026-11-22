package config

import (
	"errors"
	"time"

	"cardhub/backend/foundation/config"
)

type Config struct {
	// Application
	AppEnv   string
	GRPCAddr string

	// Database
	DatabaseURL string

	// JWT signing
	JWTKeyPath   string
	JWTIssuer    string
	JWTAudience  string
	JWTAccessTTL time.Duration

	// Refresh sessions
	JWTRefreshTTL time.Duration
}

func Load() (Config, error) {
	databaseURL, err := config.Required("DATABASE_URL")
	if err != nil {
		return Config{}, err
	}

	keyPath, err := config.Required("JWT_KEY_PATH")
	if err != nil {
		return Config{}, err
	}

	accessTTL, err := config.DurationEnv(
		"JWT_ACCESS_TTL",
		15*time.Minute,
	)
	if err != nil {
		return Config{}, err
	}

	refreshTTL, err := config.DurationEnv(
		"REFRESH_TOKEN_TTL",
		30*24*time.Hour,
	)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		AppEnv:   config.Env("APP_ENV", "development"),
		GRPCAddr: config.Env("GRPC_ADDR", ":9001"),

		DatabaseURL: databaseURL,

		JWTKeyPath:   keyPath,
		JWTIssuer:    config.Env("JWT_ISSUER", "cardhub-identity"),
		JWTAudience:  config.Env("JWT_AUDIENCE", "cardhub-api"),
		JWTAccessTTL: accessTTL,

		JWTRefreshTTL: refreshTTL,
	}

	// Validate authentication configuration at startup.
	if cfg.JWTIssuer == "" || cfg.JWTAudience == "" {
		return Config{}, errors.New(
			"JWT_ISSUER and JWT_AUDIENCE must not be empty",
		)
	}

	if cfg.JWTAccessTTL <= 0 {
		return Config{}, errors.New(
			"JWT_ACCESS_TTL must be positive",
		)
	}

	if cfg.JWTRefreshTTL <= 0 {
		return Config{}, errors.New(
			"REFRESH_TOKEN_TTL must be positive",
		)
	}

	return cfg, nil
}
