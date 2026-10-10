package locationclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type DriverLocation struct {
	DriverID  string  `json:"driver_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Client interface {
	NearbyDrivers(ctx context.Context, request NearbyDriversRequest) ([]DriverLocation, error)
	SetDriverUnavailable(ctx context.Context, driverID string) error
}

type NearbyDriversRequest struct {
	Latitude  float64
	Longitude float64
	RadiusKM  float64
	Limit     int
}

type HTTPClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *HTTPClient) NearbyDrivers(ctx context.Context, request NearbyDriversRequest) ([]DriverLocation, error) {
	endpoint, err := url.Parse(c.baseURL + "/api/v1/drivers/nearby")
	if err != nil {
		return nil, err
	}

	query := endpoint.Query()
	query.Set("latitude", strconv.FormatFloat(request.Latitude, 'f', -1, 64))
	query.Set("longitude", strconv.FormatFloat(request.Longitude, 'f', -1, 64))
	query.Set("radius_km", strconv.FormatFloat(request.RadiusKM, 'f', -1, 64))
	query.Set("limit", strconv.Itoa(request.Limit))
	endpoint.RawQuery = query.Encode()

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	response, err := c.client.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("location service returned status %d", response.StatusCode)
	}

	var drivers []DriverLocation
	if err := json.NewDecoder(response.Body).Decode(&drivers); err != nil {
		return nil, err
	}

	return drivers, nil
}

func (c *HTTPClient) SetDriverUnavailable(ctx context.Context, driverID string) error {
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/drivers/"+driverID+"/unavailable", nil)
	if err != nil {
		return err
	}
	httpRequest.Header.Set("X-User-ID", driverID)
	httpRequest.Header.Set("X-User-Role", "driver")

	response, err := c.client.Do(httpRequest)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("location service returned status %d", response.StatusCode)
	}

	return nil
}
