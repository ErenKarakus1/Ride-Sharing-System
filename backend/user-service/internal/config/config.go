package config

import "os"

type Config struct {
	Port string
}

func Load() Config {
	return Config{
		Port: env("PORT", "8080"),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
