package config

import (
	"errors"
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
	Environment     string
}

func Load() Config {
	return Config{
		HTTPPort:        env("HTTP_PORT", env("PORT", "8081")),
		GRPCPort:        env("GRPC_PORT", "9092"),
		DatabaseURL:     env("DATABASE_URL", "postgres://rideshare:rideshare@localhost:5432/rideshare?sslmode=disable"),
		JWTSecret:       env("JWT_SECRET", "change-me"),
		TokenTTL:        24 * time.Hour,
		UserServiceAddr: env("USER_SERVICE_GRPC_ADDR", "localhost:9091"),
		Environment:     env("APP_ENV", "local"),
	}
}

func (c Config) Validate() error {
	if c.Environment == "local" || c.Environment == "development" || c.Environment == "test" {
		return nil
	}
	if c.JWTSecret == "" || c.JWTSecret == "change-me" || len(c.JWTSecret) < 32 {
		return errors.New("JWT_SECRET must be set to at least 32 characters outside local development")
	}

	return nil
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
