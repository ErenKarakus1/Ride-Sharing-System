package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
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
	ID       string `json:"id"`
	DriverID string `json:"driver_id"`
	Status   string `json:"status"`
}

type paymentResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type notificationMessage struct {
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
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
		"password":     "Password12345",
		"display_name": "E2E Rider",
		"phone_number": "+905551112233",
		"role":         "rider",
	})
	must("register", err)

	login, err := post[authResponse](client, baseURL+"/api/v1/auth/login", "", map[string]any{
		"email":    email,
		"password": "Password12345",
	})
	must("login", err)

	driverEmail := fmt.Sprintf("driver%d@example.com", time.Now().UnixNano())
	driverRegister, err := post[authResponse](client, baseURL+"/api/v1/auth/register", "", map[string]any{
		"email":        driverEmail,
		"password":     "Password12345",
		"display_name": "E2E Driver",
		"phone_number": "+905559998877",
		"role":         "driver",
	})
	must("register driver", err)

	driverLogin, err := post[authResponse](client, baseURL+"/api/v1/auth/login", "", map[string]any{
		"email":    driverEmail,
		"password": "Password12345",
	})
	must("login driver", err)

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

	notifications, closeNotifications, err := watchNotifications(baseURL, login.AccessToken)
	must("connect rider notifications websocket", err)
	defer closeNotifications()

	must("update driver location", putNoContent(client, fmt.Sprintf("%s/api/v1/drivers/%s/location", baseURL, driverRegister.UserID), driverLogin.AccessToken, map[string]any{
		"latitude":  pickup.Latitude,
		"longitude": pickup.Longitude,
	}))
	must("set driver available", postNoContent(client, fmt.Sprintf("%s/api/v1/drivers/%s/available", baseURL, driverRegister.UserID), driverLogin.AccessToken, nil))

	payment, err := post[paymentResponse](client, baseURL+"/api/v1/payments/authorize", login.AccessToken, map[string]any{
		"ride_id":  ride.ID,
		"amount":   fare.Amount,
		"currency": fare.Currency,
	})
	must("authorize payment", err)

	accepted, err := waitForRideStatus(client, fmt.Sprintf("%s/api/v1/rides/%s", baseURL, ride.ID), login.AccessToken, "accepted")
	must("wait for auto-accepted ride", err)
	if accepted.DriverID != driverRegister.UserID {
		must("auto-match assigned expected driver", fmt.Errorf("expected driver %s, got %s", driverRegister.UserID, accepted.DriverID))
	}

	driver2Email := fmt.Sprintf("driver2-%d@example.com", time.Now().UnixNano())
	driver2, err := post[authResponse](client, baseURL+"/api/v1/auth/register", "", map[string]any{
		"email":        driver2Email,
		"password":     "Password12345",
		"display_name": "E2E Second Driver",
		"phone_number": "+905550000002",
		"role":         "driver",
	})
	must("register second driver", err)
	driver2Login, err := post[authResponse](client, baseURL+"/api/v1/auth/login", "", map[string]any{
		"email":    driver2Email,
		"password": "Password12345",
	})
	must("login second driver", err)
	must("second driver cannot accept already accepted ride", postExpectStatus(client, fmt.Sprintf("%s/api/v1/rides/%s/accept", baseURL, ride.ID), driver2Login.AccessToken, map[string]any{"driver_id": driver2.UserID}, http.StatusConflict))

	secondRide, err := post[rideResponse](client, baseURL+"/api/v1/rides", login.AccessToken, map[string]any{
		"pickup":  pickup,
		"dropoff": dropoff,
	})
	must("create second ride", err)
	secondPayment, err := post[paymentResponse](client, baseURL+"/api/v1/payments/authorize", login.AccessToken, map[string]any{
		"ride_id":  secondRide.ID,
		"amount":   fare.Amount,
		"currency": fare.Currency,
	})
	must("authorize second payment", err)
	_ = secondPayment
	_, err = waitForRideStatusNot(client, fmt.Sprintf("%s/api/v1/rides/%s", baseURL, secondRide.ID), login.AccessToken, "accepted", 3*time.Second)
	must("second ride stays unmatched while driver is claimed", err)
	must("claimed driver cannot match second ride", postExpectStatus(client, baseURL+"/api/v1/matches", login.AccessToken, map[string]any{
		"ride_id":   secondRide.ID,
		"pickup":    pickup,
		"radius_km": 5,
		"limit":     5,
	}, http.StatusNotFound))

	must("driver unavailable reset prevents matching", postNoContent(client, fmt.Sprintf("%s/api/v1/drivers/%s/unavailable", baseURL, driverRegister.UserID), driverLogin.AccessToken, nil))
	must("unavailable driver cannot match", postExpectStatus(client, baseURL+"/api/v1/matches", login.AccessToken, map[string]any{
		"ride_id":   secondRide.ID,
		"pickup":    pickup,
		"radius_km": 5,
		"limit":     5,
	}, http.StatusNotFound))
	_, err = post[rideResponse](client, fmt.Sprintf("%s/api/v1/rides/%s/cancel", baseURL, secondRide.ID), login.AccessToken, map[string]any{})
	must("cancel unmatched second ride", err)

	started, err := post[rideResponse](client, fmt.Sprintf("%s/api/v1/rides/%s/start", baseURL, ride.ID), driverLogin.AccessToken, map[string]any{})
	must("start ride", err)

	completed, err := post[rideResponse](client, fmt.Sprintf("%s/api/v1/rides/%s/complete", baseURL, ride.ID), driverLogin.AccessToken, map[string]any{})
	must("complete ride", err)

	must("wait for ride completed notification", waitForNotification(notifications, "ride.completed", ride.ID))

	captured, err := waitForPaymentStatus(client, fmt.Sprintf("%s/api/v1/payments/%s", baseURL, payment.ID), login.AccessToken, "captured")
	must("wait for captured payment", err)

	fmt.Printf("E2E passed: rider=%s driver=%s matched_driver=%v ride=%s statuses=%s/%s/%s/%s fare=%.2f %s payment=%s payment_status=%s/%s\n",
		register.UserID,
		driverRegister.UserID,
		accepted.DriverID,
		ride.ID,
		ride.Status,
		accepted.Status,
		started.Status,
		completed.Status,
		fare.Amount,
		fare.Currency,
		payment.ID,
		payment.Status,
		captured.Status,
	)
}

func waitForRideStatusNot(client *http.Client, url string, accessToken string, status string, duration time.Duration) (rideResponse, error) {
	deadline := time.After(duration)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var ride rideResponse
	for {
		var err error
		ride, err = get[rideResponse](client, url, accessToken)
		if err != nil {
			return rideResponse{}, err
		}
		if ride.Status == status {
			return ride, fmt.Errorf("ride unexpectedly became %q", status)
		}

		select {
		case <-deadline:
			return ride, nil
		case <-ticker.C:
		}
	}
}

func waitForRideStatus(client *http.Client, url string, accessToken string, status string) (rideResponse, error) {
	var ride rideResponse
	var err error
	deadline := time.After(60 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		ride, err = get[rideResponse](client, url, accessToken)
		if err != nil {
			return rideResponse{}, err
		}
		if ride.Status == status {
			return ride, nil
		}

		select {
		case <-deadline:
			return ride, fmt.Errorf("ride status remained %q, expected %q after 60s", ride.Status, status)
		case <-ticker.C:
		}
	}
}

func watchNotifications(baseURL string, accessToken string) (<-chan notificationMessage, func(), error) {
	wsURL := strings.Replace(baseURL, "http://", "ws://", 1)
	wsURL = strings.Replace(wsURL, "https://", "wss://", 1)
	wsURL += "/ws/notifications"

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+accessToken)
	headers.Set("Origin", "http://localhost:3000")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
	if err != nil {
		return nil, nil, err
	}

	notifications := make(chan notificationMessage, 8)
	done := make(chan struct{})
	go func() {
		defer close(notifications)
		for {
			select {
			case <-done:
				return
			default:
			}

			var message notificationMessage
			if err := conn.ReadJSON(&message); err != nil {
				return
			}
			notifications <- message
		}
	}()

	closeFn := func() {
		close(done)
		_ = conn.Close()
	}

	return notifications, closeFn, nil
}

func waitForNotification(notifications <-chan notificationMessage, eventType string, rideID string) error {
	timeout := time.After(10 * time.Second)
	for {
		select {
		case message, ok := <-notifications:
			if !ok {
				return fmt.Errorf("notifications websocket closed before %s", eventType)
			}
			if message.Type != eventType {
				continue
			}
			if fmt.Sprint(message.Payload["ride_id"]) == rideID {
				return nil
			}
		case <-timeout:
			return fmt.Errorf("timed out waiting for %s notification", eventType)
		}
	}
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

func get[T any](client *http.Client, url string, accessToken string) (T, error) {
	var result T

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return result, err
	}
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

func waitForPaymentStatus(client *http.Client, url string, accessToken string, status string) (paymentResponse, error) {
	var payment paymentResponse
	var err error
	deadline := time.After(60 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		payment, err = get[paymentResponse](client, url, accessToken)
		if err != nil {
			return paymentResponse{}, err
		}
		if payment.Status == status {
			return payment, nil
		}

		select {
		case <-deadline:
			return payment, fmt.Errorf("payment status remained %q, expected %q after 60s", payment.Status, status)
		case <-ticker.C:
		}
	}
}

func postNoContent(client *http.Client, url string, accessToken string, body any) error {
	return requestNoContent(client, http.MethodPost, url, accessToken, body)
}

func postExpectStatus(client *http.Client, url string, accessToken string, body any, expectedStatus int) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+accessToken)

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != expectedStatus {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("%s returned status %d, expected %d: %s", url, response.StatusCode, expectedStatus, strings.TrimSpace(string(body)))
	}

	return nil
}

func putNoContent(client *http.Client, url string, accessToken string, body any) error {
	return requestNoContent(client, http.MethodPut, url, accessToken, body)
}

func requestNoContent(client *http.Client, method string, url string, accessToken string, body any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(method, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+accessToken)

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return fmt.Errorf("%s returned status %d", url, response.StatusCode)
	}

	return nil
}

func must(step string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s failed: %v\n", step, err)
		os.Exit(1)
	}
}
