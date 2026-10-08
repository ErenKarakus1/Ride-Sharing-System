package config

import "os"

type Config struct {
	Port          string
	KafkaBrokers  string
	ConsumerGroup string
}

func Load() Config {
	return Config{
		Port:          env("PORT", "8085"),
		KafkaBrokers:  env("KAFKA_BROKERS", "localhost:9092"),
		ConsumerGroup: env("KAFKA_CONSUMER_GROUP", "notification-service"),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
