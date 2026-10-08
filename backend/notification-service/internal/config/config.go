package config

import (
	"os"
	"strings"
)

type Config struct {
	Port                string
	KafkaBrokers        string
	ConsumerGroup       string
	AuthServiceGRPCAddr string
	AllowedOrigins      []string
}

func Load() Config {
	return Config{
		Port:                env("PORT", "8085"),
		KafkaBrokers:        env("KAFKA_BROKERS", "localhost:9092"),
		ConsumerGroup:       env("KAFKA_CONSUMER_GROUP", "notification-service"),
		AuthServiceGRPCAddr: env("AUTH_SERVICE_GRPC_ADDR", "localhost:9092"),
		AllowedOrigins:      envList("WS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173"),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func envList(key string, fallback string) []string {
	value := env(key, fallback)
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			items = append(items, part)
		}
	}

	return items
}
