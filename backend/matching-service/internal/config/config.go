package config

import "os"

type Config struct {
	Port               string
	LocationServiceURL string
}

func Load() Config {
	return Config{
		Port:               env("PORT", "8084"),
		LocationServiceURL: env("LOCATION_SERVICE_HTTP_URL", "http://localhost:8083"),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
