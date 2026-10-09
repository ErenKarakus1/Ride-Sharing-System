package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type authResponse struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"user_id"`
}

type location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address,omitempty"`
}

type fareEstimate struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type rideResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type paymentResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func main() {
	baseURL := os.Getenv("API_GATEWAY_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8088"
	}

	client := &http.Client{Timeout: 10 * time.Second}
	email := fmt.Sprintf("rider%d@example.com", time.Now().UnixNano())

	register, err := post[authResponse](client, baseURL+"/api/v1/auth/register", "", map[string]any{
		"email":        email,
		"password":     "password123",
		"display_name": "E2E Rider",
		"phone_number": "+905551112233",
		"role":         "rider",
	})
	must("register", err)

	login, err := post[authResponse](client, baseURL+"/api/v1/auth/login", "", map[string]any{
		"email":    email,
		"password": "password123",
	})
	must("login", err)

	pickup := location{Latitude: 41.0082, Longitude: 28.9784, Address: "Sultanahmet"}
	dropoff := location{Latitude: 41.0369, Longitude: 28.9850, Address: "Taksim"}

	fare, err := post[fareEstimate](client, baseURL+"/api/v1/fare-estimates", login.AccessToken, map[string]any{
		"pickup":  pickup,
		"dropoff": dropoff,
	})
	must("fare estimate", err)

	ride, err := post[rideResponse](client, baseURL+"/api/v1/rides", login.AccessToken, map[string]any{
		"pickup":  pickup,
		"dropoff": dropoff,
	})
	must("create ride", err)

	payment, err := post[paymentResponse](client, baseURL+"/api/v1/payments/authorize", login.AccessToken, map[string]any{
		"ride_id":  ride.ID,
		"rider_id": register.UserID,
		"amount":   fare.Amount,
		"currency": fare.Currency,
	})
	must("authorize payment", err)

	fmt.Printf("E2E passed: user=%s ride=%s ride_status=%s fare=%.2f %s payment=%s payment_status=%s\n",
		register.UserID,
		ride.ID,
		ride.Status,
		fare.Amount,
		fare.Currency,
		payment.ID,
		payment.Status,
	)
}

func post[T any](client *http.Client, url string, accessToken string, body any) (T, error) {
	var result T

	payload, err := json.Marshal(body)
	if err != nil {
		return result, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return result, err
	}
	request.Header.Set("Content-Type", "application/json")
	if accessToken != "" {
		request.Header.Set("Authorization", "Bearer "+accessToken)
	}

	response, err := client.Do(request)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return result, fmt.Errorf("%s returned status %d", url, response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}

func must(step string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s failed: %v\n", step, err)
		os.Exit(1)
	}
}
