package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port      string
	BaseFare  float64
	PerKM     float64
	PerMinute float64
}

func Load() Config {
	return Config{
		Port:      env("PORT", "8086"),
		BaseFare:  envFloat("BASE_FARE", 35),
		PerKM:     envFloat("FARE_PER_KM", 12),
		PerMinute: envFloat("FARE_PER_MINUTE", 2),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func envFloat(key string, fallback float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}

	return parsed
}
