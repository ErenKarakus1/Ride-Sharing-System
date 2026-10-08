package matching

import (
	"context"
	"errors"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/matching-service/internal/locationclient"
)

var ErrNoDriversAvailable = errors.New("no drivers available")

type Service struct {
	locations locationclient.Client
}

func NewService(locations locationclient.Client) *Service {
	return &Service{locations: locations}
}

func (s *Service) Match(ctx context.Context, request MatchRequest) (MatchResponse, error) {
	limit := request.Limit
	if limit <= 0 {
		limit = 5
	}

	radiusKM := request.RadiusKM
	if radiusKM <= 0 {
		radiusKM = 5
	}

	drivers, err := s.locations.NearbyDrivers(ctx, locationclient.NearbyDriversRequest{
		Latitude:  request.Pickup.Latitude,
		Longitude: request.Pickup.Longitude,
		RadiusKM:  radiusKM,
		Limit:     limit,
	})
	if err != nil {
		return MatchResponse{}, err
	}
	if len(drivers) == 0 {
		return MatchResponse{}, ErrNoDriversAvailable
	}

	driver := drivers[0]
	return MatchResponse{
		RideID:    request.RideID,
		DriverID:  driver.DriverID,
		Latitude:  driver.Latitude,
		Longitude: driver.Longitude,
	}, nil
}
