package location

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	claimDriverID       string
	claimErr            error
	nearbyLatitude      float64
	nearbyLongitude     float64
	nearbyRadiusKM      float64
	nearbyLimit         int
	nearbyDrivers       []DriverLocation
	nearbyErr           error
	availableDriverID   string
	unavailableDriverID string
}

func (f *fakeRepository) UpdateDriverLocation(ctx context.Context, location DriverLocation) error {
	return nil
}

func (f *fakeRepository) SetDriverAvailable(ctx context.Context, driverID string) error {
	f.availableDriverID = driverID
	return nil
}

func (f *fakeRepository) SetDriverUnavailable(ctx context.Context, driverID string) error {
	f.unavailableDriverID = driverID
	return nil
}

func (f *fakeRepository) ClaimDriver(ctx context.Context, driverID string) error {
	f.claimDriverID = driverID
	return f.claimErr
}

func (f *fakeRepository) NearbyDrivers(ctx context.Context, latitude float64, longitude float64, radiusKM float64, limit int) ([]DriverLocation, error) {
	f.nearbyLatitude = latitude
	f.nearbyLongitude = longitude
	f.nearbyRadiusKM = radiusKM
	f.nearbyLimit = limit
	return f.nearbyDrivers, f.nearbyErr
}

func TestClaimDriverTrimsDriverID(t *testing.T) {
	repository := &fakeRepository{}
	service := NewService(repository)

	if err := service.ClaimDriver(context.Background(), " driver-1 "); err != nil {
		t.Fatalf("ClaimDriver returned error: %v", err)
	}
	if repository.claimDriverID != "driver-1" {
		t.Fatalf("expected trimmed driver ID, got %q", repository.claimDriverID)
	}
}

func TestClaimDriverRequiresDriverID(t *testing.T) {
	service := NewService(&fakeRepository{})

	err := service.ClaimDriver(context.Background(), "   ")
	if !errors.Is(err, ErrMissingDriver) {
		t.Fatalf("expected ErrMissingDriver, got %v", err)
	}
}

func TestClaimDriverPropagatesUnavailable(t *testing.T) {
	service := NewService(&fakeRepository{claimErr: ErrDriverUnavailable})

	err := service.ClaimDriver(context.Background(), "driver-1")
	if !errors.Is(err, ErrDriverUnavailable) {
		t.Fatalf("expected ErrDriverUnavailable, got %v", err)
	}
}

func TestNearbyDriversDefaultsSearchOptions(t *testing.T) {
	repository := &fakeRepository{
		nearbyDrivers: []DriverLocation{
			{DriverID: "driver-1", Latitude: 41.0082, Longitude: 28.9784},
		},
	}
	service := NewService(repository)

	drivers, err := service.NearbyDrivers(context.Background(), NearbyDriversRequest{
		Latitude:  41.0082,
		Longitude: 28.9784,
		RadiusKM:  -1,
		Limit:     100,
	})
	if err != nil {
		t.Fatalf("NearbyDrivers returned error: %v", err)
	}
	if len(drivers) != 1 || drivers[0].DriverID != "driver-1" {
		t.Fatalf("unexpected drivers: %#v", drivers)
	}
	if repository.nearbyRadiusKM != 5 {
		t.Fatalf("expected default radius 5km, got %v", repository.nearbyRadiusKM)
	}
	if repository.nearbyLimit != 10 {
		t.Fatalf("expected default limit 10, got %d", repository.nearbyLimit)
	}
}

func TestAvailabilityMethodsRequireDriverID(t *testing.T) {
	service := NewService(&fakeRepository{})

	if err := service.SetDriverAvailable(context.Background(), " "); !errors.Is(err, ErrMissingDriver) {
		t.Fatalf("SetDriverAvailable expected ErrMissingDriver, got %v", err)
	}
	if err := service.SetDriverUnavailable(context.Background(), " "); !errors.Is(err, ErrMissingDriver) {
		t.Fatalf("SetDriverUnavailable expected ErrMissingDriver, got %v", err)
	}
}
