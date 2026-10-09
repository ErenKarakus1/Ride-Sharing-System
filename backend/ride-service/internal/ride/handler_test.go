package ride

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/ride-service/internal/events"
	"github.com/gin-gonic/gin"
)

func TestHandlerRequiresDriverRoleToAcceptRide(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := NewService(newFakeRepository(), noopPublisher{})
	handler := NewHandler(service)

	created, err := service.Create(context.Background(), CreateRideRequest{
		RiderID: "rider-1",
		Pickup:  Location{Latitude: 41, Longitude: 29},
		Dropoff: Location{Latitude: 42, Longitude: 30},
	})
	if err != nil {
		t.Fatalf("create ride: %v", err)
	}

	router := gin.New()
	router.POST("/rides/:id/accept", handler.Accept)

	request := httptest.NewRequest(http.MethodPost, "/rides/"+created.ID+"/accept", bytes.NewReader([]byte(`{}`)))
	request.Header.Set("X-User-ID", "rider-1")
	request.Header.Set("X-User-Role", "rider")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden, got %d", response.Code)
	}
}

func TestHandlerAcceptsRideAsDriver(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := NewService(newFakeRepository(), noopPublisher{})
	handler := NewHandler(service)

	created, err := service.Create(context.Background(), CreateRideRequest{
		RiderID: "rider-1",
		Pickup:  Location{Latitude: 41, Longitude: 29},
		Dropoff: Location{Latitude: 42, Longitude: 30},
	})
	if err != nil {
		t.Fatalf("create ride: %v", err)
	}

	router := gin.New()
	router.POST("/rides/:id/accept", handler.Accept)

	request := httptest.NewRequest(http.MethodPost, "/rides/"+created.ID+"/accept", bytes.NewReader([]byte(`{}`)))
	request.Header.Set("X-User-ID", "driver-1")
	request.Header.Set("X-User-Role", "driver")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected ok, got %d body=%s", response.Code, response.Body.String())
	}

	var accepted Ride
	if err := json.NewDecoder(response.Body).Decode(&accepted); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if accepted.DriverID == nil || *accepted.DriverID != "driver-1" {
		t.Fatalf("expected driver id from header, got %#v", accepted.DriverID)
	}
}

var _ events.Publisher = noopPublisher{}
