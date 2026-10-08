package config

import "os"

type Config struct {
	Port                   string
	AuthServiceGRPCAddr    string
	AuthServiceHTTPURL     string
	UserServiceHTTPURL     string
	RideServiceHTTPURL     string
	LocationServiceHTTPURL string
}

func Load() Config {
	return Config{
		Port:                   env("PORT", "8088"),
		AuthServiceGRPCAddr:    env("AUTH_SERVICE_GRPC_ADDR", "localhost:9092"),
		AuthServiceHTTPURL:     env("AUTH_SERVICE_HTTP_URL", "http://localhost:8081"),
		UserServiceHTTPURL:     env("USER_SERVICE_HTTP_URL", "http://localhost:8080"),
		RideServiceHTTPURL:     env("RIDE_SERVICE_HTTP_URL", "http://localhost:8082"),
		LocationServiceHTTPURL: env("LOCATION_SERVICE_HTTP_URL", "http://localhost:8083"),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
