package matching

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/matching-service/internal/events"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/matching-service/internal/locationclient"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/matching-service/internal/rideclient"
)

var ErrNoDriversAvailable = errors.New("no drivers available")

type Service struct {
	locations     locationclient.Client
	rides         rideclient.Client
	activeMatches map[string]struct{}
	mu            sync.Mutex
}

func NewService(locations locationclient.Client, rides rideclient.Client) *Service {
	return &Service{
		locations:     locations,
		rides:         rides,
		activeMatches: make(map[string]struct{}),
	}
}

func (s *Service) Match(ctx context.Context, request MatchRequest) (MatchResponse, error) {
	match, _, err := s.match(ctx, request, false)
	return match, err
}

func (s *Service) MatchAndAssign(ctx context.Context, request MatchRequest) (MatchResponse, error) {
	match, _, err := s.match(ctx, request, true)
	return match, err
}

func (s *Service) match(ctx context.Context, request MatchRequest, assign bool) (MatchResponse, rideclient.Ride, error) {
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
		return MatchResponse{}, rideclient.Ride{}, err
	}
	if len(drivers) == 0 {
		return MatchResponse{}, rideclient.Ride{}, ErrNoDriversAvailable
	}

	driver := drivers[0]
	match := MatchResponse{
		RideID:    request.RideID,
		DriverID:  driver.DriverID,
		Latitude:  driver.Latitude,
		Longitude: driver.Longitude,
	}

	if !assign || s.rides == nil {
		return match, rideclient.Ride{}, nil
	}

	accepted, err := s.rides.AcceptRide(ctx, request.RideID, driver.DriverID)
	if err != nil {
		return MatchResponse{}, rideclient.Ride{}, err
	}
	if err := s.locations.SetDriverUnavailable(ctx, driver.DriverID); err != nil {
		return MatchResponse{}, rideclient.Ride{}, err
	}

	return match, accepted, nil
}

func (s *Service) HandlePaymentEvent(ctx context.Context, event events.PaymentEvent) error {
	if event.Type != "payment.authorized" {
		return nil
	}
	if s.rides == nil {
		log.Printf("skipping background match for ride_id=%s: ride client is not configured", event.Data.RideID)
		return nil
	}

	ride, err := s.rides.GetRide(ctx, event.Data.RideID)
	if err != nil {
		return err
	}
	if ride.Status != "requested" || ride.DriverID != nil {
		log.Printf("skipping background match for ride_id=%s: status=%s driver_assigned=%t", ride.ID, ride.Status, ride.DriverID != nil)
		return nil
	}

	if !s.startBackgroundMatch(ride.ID) {
		log.Printf("background match already running for ride_id=%s", ride.ID)
		return nil
	}

	request := MatchRequest{
		RideID:   ride.ID,
		Pickup:   Location{Latitude: ride.Pickup.Latitude, Longitude: ride.Pickup.Longitude},
		RadiusKM: 5,
		Limit:    5,
	}

	log.Printf("starting background match for ride_id=%s", ride.ID)
	go s.runBackgroundMatch(request)
	return nil
}

func (s *Service) runBackgroundMatch(request MatchRequest) {
	defer s.finishBackgroundMatch(request.RideID)

	match, accepted, err := s.matchWithRetry(context.Background(), request)
	if err != nil {
		log.Printf("background match failed for ride_id=%s: %v", request.RideID, err)
		return
	}
	if match.DriverID == "" {
		return
	}

	log.Printf("auto-matched ride_id=%s driver_id=%s status=%s", accepted.ID, match.DriverID, accepted.Status)
}

func (s *Service) matchWithRetry(ctx context.Context, request MatchRequest) (MatchResponse, rideclient.Ride, error) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		ride, err := s.rides.GetRide(ctx, request.RideID)
		if err != nil {
			return MatchResponse{}, rideclient.Ride{}, err
		}
		if ride.Status != "requested" || ride.DriverID != nil {
			log.Printf("stopping background match for ride_id=%s: status=%s driver_assigned=%t", ride.ID, ride.Status, ride.DriverID != nil)
			return MatchResponse{}, ride, nil
		}

		match, accepted, err := s.match(ctx, request, true)
		if err == nil {
			return match, accepted, nil
		}
		if !errors.Is(err, ErrNoDriversAvailable) {
			return MatchResponse{}, rideclient.Ride{}, err
		}

		log.Printf("no drivers available for ride_id=%s; retrying background match", request.RideID)
		select {
		case <-ctx.Done():
			return MatchResponse{}, rideclient.Ride{}, nil
		case <-ticker.C:
		}
	}
}

func (s *Service) startBackgroundMatch(rideID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.activeMatches[rideID]; exists {
		return false
	}

	s.activeMatches[rideID] = struct{}{}
	return true
}

func (s *Service) finishBackgroundMatch(rideID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.activeMatches, rideID)
}
