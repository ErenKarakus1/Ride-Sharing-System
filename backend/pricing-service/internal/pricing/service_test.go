package pricing

import "testing"

func TestEstimate(t *testing.T) {
	service := NewService(35, 12, 2)

	estimate := service.Estimate(EstimateRequest{
		Pickup:  Location{Latitude: 41.0082, Longitude: 28.9784},
		Dropoff: Location{Latitude: 41.0369, Longitude: 28.9850},
	})

	if estimate.Currency != "TRY" {
		t.Fatalf("expected TRY, got %s", estimate.Currency)
	}
	if estimate.DistanceKM <= 0 {
		t.Fatalf("expected positive distance, got %f", estimate.DistanceKM)
	}
	if estimate.Amount <= 35 {
		t.Fatalf("expected amount above base fare, got %f", estimate.Amount)
	}
}
