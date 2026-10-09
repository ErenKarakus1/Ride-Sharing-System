package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type healthResponse struct {
	Status string `json:"status"`
}

type service struct {
	name string
	url  string
}

func main() {
	apiGatewayURL := os.Getenv("API_GATEWAY_URL")
	if apiGatewayURL == "" {
		apiGatewayURL = "http://localhost:8088"
	}

	services := []service{
		{name: "api-gateway", url: apiGatewayURL + "/health"},
		{name: "user-service", url: "http://localhost:8080/health"},
		{name: "auth-service", url: "http://localhost:8081/health"},
		{name: "ride-service", url: "http://localhost:8082/health"},
		{name: "location-service", url: "http://localhost:8083/health"},
		{name: "matching-service", url: "http://localhost:8084/health"},
		{name: "notification-service", url: "http://localhost:8085/health"},
		{name: "pricing-service", url: "http://localhost:8086/health"},
		{name: "payment-service", url: "http://localhost:8087/health"},
	}

	client := &http.Client{Timeout: 10 * time.Second}
	for _, service := range services {
		fmt.Printf("Checking %s at %s\n", service.name, service.url)
		if err := waitForHealth(client, service.url, 90*time.Second); err != nil {
			fmt.Fprintf(os.Stderr, "%s health check failed: %v\n", service.name, err)
			os.Exit(1)
		}
	}

	fmt.Println("Smoke checks passed.")
}

func waitForHealth(client *http.Client, url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		if err := checkHealth(client, url); err != nil {
			lastErr = err
			time.Sleep(2 * time.Second)
			continue
		}

		return nil
	}

	return lastErr
}

func checkHealth(client *http.Client, url string) error {
	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", response.StatusCode)
	}

	var health healthResponse
	if err := json.NewDecoder(response.Body).Decode(&health); err != nil {
		return err
	}
	if health.Status != "ok" {
		return fmt.Errorf("unexpected health status %q", health.Status)
	}

	return nil
}
