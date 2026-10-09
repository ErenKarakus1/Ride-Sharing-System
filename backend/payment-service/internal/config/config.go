package config

import "os"

type Config struct {
	Port                string
	DatabaseURL         string
	KafkaBrokers        string
	ConsumerGroup       string
	RideServiceGRPCAddr string
	InternalToken       string
}

func Load() Config {
	return Config{
		Port:                env("PORT", "8087"),
		DatabaseURL:         env("DATABASE_URL", "postgres://rideshare:rideshare@localhost:5432/rideshare?sslmode=disable"),
		KafkaBrokers:        env("KAFKA_BROKERS", "localhost:9092"),
		ConsumerGroup:       env("KAFKA_CONSUMER_GROUP", "payment-service"),
		RideServiceGRPCAddr: env("RIDE_SERVICE_GRPC_ADDR", "localhost:9093"),
		InternalToken:       env("INTERNAL_SERVICE_TOKEN", ""),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
