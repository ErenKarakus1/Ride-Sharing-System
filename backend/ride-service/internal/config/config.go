package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() Config {
	return Config{
		Port:        env("PORT", "8082"),
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
