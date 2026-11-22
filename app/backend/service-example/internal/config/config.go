package config

import (
	foundationconfig "cardhub/backend/foundation/config"
)

type Config struct {
	// Application
	AppEnv   string
	GRPCAddr string

	// Database
	DatabaseURL string
}

func Load() (Config, error) {
	databaseURL, err := foundationconfig.Required("DATABASE_URL")
	if err != nil {
		return Config{}, err
	}

	return Config{
		AppEnv:      foundationconfig.Env("APP_ENV", "development"),
		GRPCAddr:    foundationconfig.Env("GRPC_ADDR", ":9000"),
		DatabaseURL: databaseURL,
	}, nil
}
