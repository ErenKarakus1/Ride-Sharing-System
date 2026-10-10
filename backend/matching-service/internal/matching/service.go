package matching

import (
	"context"
	"errors"
	"log"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/matching-service/internal/events"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/matching-service/internal/locationclient"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/matching-service/internal/rideclient"
)

var ErrNoDriversAvailable = errors.New("no drivers available")

type Service struct {
	locations locationclient.Client
	rides     rideclient.Client
}

func NewService(locations locationclient.Client, rides rideclient.Client) *Service {
	return &Service{locations: locations, rides: rides}
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

func (s *Service) HandlePaymentEvent(ctx context.Context, event events.PaymentEvent) error {
	if event.Type != "payment.authorized" {
		return nil
	}
	if s.rides == nil {
		return nil
	}

	ride, err := s.rides.GetRide(ctx, event.Data.RideID)
	if err != nil {
		return err
	}
	if ride.Status != "requested" || ride.DriverID != nil {
		return nil
	}

	match, err := s.Match(ctx, MatchRequest{
		RideID:   ride.ID,
		Pickup:   Location{Latitude: ride.Pickup.Latitude, Longitude: ride.Pickup.Longitude},
		RadiusKM: 5,
		Limit:    5,
	})
	if errors.Is(err, ErrNoDriversAvailable) {
		log.Printf("no drivers available for ride_id=%s", ride.ID)
		return nil
	}
	if err != nil {
		return err
	}

	accepted, err := s.rides.AcceptRide(ctx, ride.ID, match.DriverID)
	if err != nil {
		return err
	}

	log.Printf("auto-matched ride_id=%s driver_id=%s status=%s", accepted.ID, match.DriverID, accepted.Status)
	return nil
}
