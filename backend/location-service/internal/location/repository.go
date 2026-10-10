package location

import "context"

type Repository interface {
	UpdateDriverLocation(ctx context.Context, location DriverLocation) error
	SetDriverAvailable(ctx context.Context, driverID string) error
	SetDriverUnavailable(ctx context.Context, driverID string) error
	ClaimDriver(ctx context.Context, driverID string) error
	NearbyDrivers(ctx context.Context, latitude float64, longitude float64, radiusKM float64, limit int) ([]DriverLocation, error)
}
