package config

import (
	"os"
	"time"
)

type Config struct {
	Port            string
	DatabaseURL     string
	JWTSecret       string
	TokenTTL        time.Duration
	UserServiceAddr string
}

func Load() Config {
	return Config{
		Port:            env("PORT", "8081"),
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
