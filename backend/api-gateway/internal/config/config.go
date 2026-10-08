package config

import "os"

type Config struct {
	Port                       string
	AuthServiceGRPCAddr        string
	AuthServiceHTTPURL         string
	UserServiceHTTPURL         string
	RideServiceHTTPURL         string
	LocationServiceHTTPURL     string
	MatchingServiceHTTPURL     string
	PricingServiceHTTPURL      string
	NotificationServiceHTTPURL string
	PaymentServiceHTTPURL      string
}

func Load() Config {
	return Config{
		Port:                       env("PORT", "8088"),
		AuthServiceGRPCAddr:        env("AUTH_SERVICE_GRPC_ADDR", "localhost:9092"),
		AuthServiceHTTPURL:         env("AUTH_SERVICE_HTTP_URL", "http://localhost:8081"),
		UserServiceHTTPURL:         env("USER_SERVICE_HTTP_URL", "http://localhost:8080"),
		RideServiceHTTPURL:         env("RIDE_SERVICE_HTTP_URL", "http://localhost:8082"),
		LocationServiceHTTPURL:     env("LOCATION_SERVICE_HTTP_URL", "http://localhost:8083"),
		MatchingServiceHTTPURL:     env("MATCHING_SERVICE_HTTP_URL", "http://localhost:8084"),
		PricingServiceHTTPURL:      env("PRICING_SERVICE_HTTP_URL", "http://localhost:8086"),
		NotificationServiceHTTPURL: env("NOTIFICATION_SERVICE_HTTP_URL", "http://localhost:8085"),
		PaymentServiceHTTPURL:      env("PAYMENT_SERVICE_HTTP_URL", "http://localhost:8087"),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
