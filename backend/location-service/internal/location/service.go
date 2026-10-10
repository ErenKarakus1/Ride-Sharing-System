package location

import (
	"context"
	"errors"
	"strings"
)

var ErrMissingDriver = errors.New("missing driver")
var ErrDriverUnavailable = errors.New("driver unavailable")

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) UpdateDriverLocation(ctx context.Context, driverID string, request UpdateDriverLocationRequest) error {
	driverID = strings.TrimSpace(driverID)
	if driverID == "" {
		return ErrMissingDriver
	}

	return s.repository.UpdateDriverLocation(ctx, DriverLocation{
		DriverID:  driverID,
		Latitude:  request.Latitude,
		Longitude: request.Longitude,
	})
}

func (s *Service) SetDriverAvailable(ctx context.Context, driverID string) error {
	driverID = strings.TrimSpace(driverID)
	if driverID == "" {
		return ErrMissingDriver
	}

	return s.repository.SetDriverAvailable(ctx, driverID)
}

func (s *Service) SetDriverUnavailable(ctx context.Context, driverID string) error {
	driverID = strings.TrimSpace(driverID)
	if driverID == "" {
		return ErrMissingDriver
	}

	return s.repository.SetDriverUnavailable(ctx, driverID)
}

func (s *Service) ClaimDriver(ctx context.Context, driverID string) error {
	driverID = strings.TrimSpace(driverID)
	if driverID == "" {
		return ErrMissingDriver
	}

	return s.repository.ClaimDriver(ctx, driverID)
}

func (s *Service) NearbyDrivers(ctx context.Context, request NearbyDriversRequest) ([]DriverLocation, error) {
	limit := request.Limit
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	radiusKM := request.RadiusKM
	if radiusKM <= 0 {
		radiusKM = 5
	}

	return s.repository.NearbyDrivers(ctx, request.Latitude, request.Longitude, radiusKM, limit)
}
