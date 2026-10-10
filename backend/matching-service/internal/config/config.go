package config

import "os"

type Config struct {
	Port               string
	LocationServiceURL string
	RideServiceURL     string
	KafkaBrokers       string
	ConsumerGroup      string
	InternalToken      string
}

func Load() Config {
	return Config{
		Port:               env("PORT", "8084"),
		LocationServiceURL: env("LOCATION_SERVICE_HTTP_URL", "http://localhost:8083"),
		RideServiceURL:     env("RIDE_SERVICE_HTTP_URL", "http://localhost:8082"),
		KafkaBrokers:       env("KAFKA_BROKERS", "localhost:9092"),
		ConsumerGroup:      env("KAFKA_CONSUMER_GROUP", "matching-service"),
		InternalToken:      env("INTERNAL_SERVICE_TOKEN", ""),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
