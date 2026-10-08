package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPPort        string
	GRPCPort        string
	DatabaseURL     string
	JWTSecret       string
	TokenTTL        time.Duration
	UserServiceAddr string
}

func Load() Config {
	return Config{
		HTTPPort:        env("HTTP_PORT", env("PORT", "8081")),
		GRPCPort:        env("GRPC_PORT", "9092"),
		DatabaseURL:     env("DATABASE_URL", "postgres://rideshare:rideshare@localhost:5432/rideshare?sslmode=disable"),
		JWTSecret:       env("JWT_SECRET", "change-me"),
		TokenTTL:        24 * time.Hour,
		UserServiceAddr: env("USER_SERVICE_GRPC_ADDR", "localhost:9091"),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
