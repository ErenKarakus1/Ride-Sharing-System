package config

import "os"

type Config struct {
	HTTPPort     string
	GRPCPort     string
	DatabaseURL  string
	KafkaBrokers string
}

func Load() Config {
	return Config{
		HTTPPort:     env("HTTP_PORT", env("PORT", "8082")),
		GRPCPort:     env("GRPC_PORT", "9093"),
		DatabaseURL:  env("DATABASE_URL", "postgres://rideshare:rideshare@localhost:5432/rideshare?sslmode=disable"),
		KafkaBrokers: env("KAFKA_BROKERS", "localhost:9092"),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
