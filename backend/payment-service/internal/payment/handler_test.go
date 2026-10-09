package payment

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/payment-service/internal/events"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/payment-service/internal/rideclient"
	"github.com/gin-gonic/gin"
)

func TestAuthorizeRequiresRiderRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(NewService(newFakeRepository(), noopPublisher{}, fakeRideClient{riderID: "rider-1", status: "requested"}))

	router := gin.New()
	router.POST("/payments/authorize", handler.Authorize)

	request := httptest.NewRequest(http.MethodPost, "/payments/authorize", bytes.NewReader([]byte(`{"ride_id":"ride-1","amount":10}`)))
	request.Header.Set("X-User-ID", "driver-1")
	request.Header.Set("X-User-Role", "driver")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden, got %d", response.Code)
	}
}

func TestAuthorizeUsesAuthenticatedRider(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(NewService(newFakeRepository(), noopPublisher{}, fakeRideClient{riderID: "rider-1", status: "requested"}))

	router := gin.New()
	router.POST("/payments/authorize", handler.Authorize)

	request := httptest.NewRequest(http.MethodPost, "/payments/authorize", bytes.NewReader([]byte(`{"ride_id":"ride-1","rider_id":"spoofed","amount":10}`)))
	request.Header.Set("X-User-ID", "rider-1")
	request.Header.Set("X-User-Role", "rider")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected created, got %d body=%s", response.Code, response.Body.String())
	}

	var authorized Payment
	if err := json.NewDecoder(response.Body).Decode(&authorized); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if authorized.RiderID != "rider-1" {
		t.Fatalf("expected authenticated rider, got %s", authorized.RiderID)
	}
}

var _ events.Publisher = noopPublisher{}
var _ rideclient.Client = fakeRideClient{}
