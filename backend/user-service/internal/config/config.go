package config

import "os"

type Config struct {
	HTTPPort    string
	GRPCPort    string
	DatabaseURL string
}

func Load() Config {
	return Config{
		HTTPPort:    env("HTTP_PORT", env("PORT", "8080")),
		GRPCPort:    env("GRPC_PORT", "9091"),
		DatabaseURL: env("DATABASE_URL", "postgres://rideshare:rideshare@localhost:5432/rideshare?sslmode=disable"),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
