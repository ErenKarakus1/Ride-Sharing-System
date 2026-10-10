package rideclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client interface {
	GetRide(ctx context.Context, id string) (Ride, error)
	AcceptRide(ctx context.Context, id string, driverID string) (Ride, error)
}

type HTTPClient struct {
	baseURL       string
	internalToken string
	httpClient    *http.Client
}

func NewHTTPClient(baseURL string, internalToken string) *HTTPClient {
	return &HTTPClient{
		baseURL:       strings.TrimRight(baseURL, "/"),
		internalToken: strings.TrimSpace(internalToken),
		httpClient:    &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *HTTPClient) GetRide(ctx context.Context, id string) (Ride, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/v1/rides/"+strings.TrimSpace(id), nil)
	if err != nil {
		return Ride{}, err
	}

	var ride Ride
	if err := c.doJSON(request, &ride); err != nil {
		return Ride{}, err
	}

	return ride, nil
}

func (c *HTTPClient) AcceptRide(ctx context.Context, id string, driverID string) (Ride, error) {
	payload, err := json.Marshal(map[string]string{"driver_id": strings.TrimSpace(driverID)})
	if err != nil {
		return Ride{}, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/rides/"+strings.TrimSpace(id)+"/accept", bytes.NewReader(payload))
	if err != nil {
		return Ride{}, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-User-ID", strings.TrimSpace(driverID))
	request.Header.Set("X-User-Role", "driver")
	if c.internalToken != "" {
		request.Header.Set("X-Internal-Service-Token", c.internalToken)
	}

	var ride Ride
	if err := c.doJSON(request, &ride); err != nil {
		return Ride{}, err
	}

	return ride, nil
}

func (c *HTTPClient) doJSON(request *http.Request, target any) error {
	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("ride-service returned status %d", response.StatusCode)
	}

	return json.NewDecoder(response.Body).Decode(target)
}

type Ride struct {
	ID       string   `json:"id"`
	RiderID  string   `json:"rider_id"`
	DriverID *string  `json:"driver_id,omitempty"`
	Pickup   Location `json:"pickup"`
	Dropoff  Location `json:"dropoff"`
	Status   string   `json:"status"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
