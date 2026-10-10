package ride

import (
	"context"
	"errors"
	"math"
	"strings"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/ride-service/internal/events"
)

var ErrInvalidTransition = errors.New("invalid ride status transition")
var ErrMissingRider = errors.New("missing rider")
var ErrUnauthorizedRideAction = errors.New("unauthorized ride action")
var ErrInvalidLocation = errors.New("invalid location")

type Service struct {
	repository Repository
	publisher  events.Publisher
}

func NewService(repository Repository, publisher events.Publisher) *Service {
	return &Service{
		repository: repository,
		publisher:  publisher,
	}
}

func (s *Service) Create(ctx context.Context, request CreateRideRequest) (Ride, error) {
	riderID := strings.TrimSpace(request.RiderID)
	if riderID == "" {
		return Ride{}, ErrMissingRider
	}
	if !validLocation(request.Pickup) || !validLocation(request.Dropoff) {
		return Ride{}, ErrInvalidLocation
	}

	created, err := s.repository.Create(ctx, Ride{
		RiderID: riderID,
		Pickup: Location{
			Latitude:  request.Pickup.Latitude,
			Longitude: request.Pickup.Longitude,
			Address:   strings.TrimSpace(request.Pickup.Address),
		},
		Dropoff: Location{
			Latitude:  request.Dropoff.Latitude,
			Longitude: request.Dropoff.Longitude,
			Address:   strings.TrimSpace(request.Dropoff.Address),
		},
	})
	if err != nil {
		return Ride{}, err
	}

	if err := s.publishRideEvent(ctx, "ride.requested", created); err != nil {
		return Ride{}, err
	}

	return created, nil
}

func (s *Service) Get(ctx context.Context, id string) (Ride, error) {
	return s.repository.Get(ctx, strings.TrimSpace(id))
}

func (s *Service) ListByRider(ctx context.Context, riderID string) ([]Ride, error) {
	return s.repository.ListByRider(ctx, strings.TrimSpace(riderID))
}

func (s *Service) Accept(ctx context.Context, id string, driverID string) (Ride, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return Ride{}, err
	}
	if current.Status != StatusRequested {
		if current.Status == StatusAccepted && driverOwnsRide(current, driverID) {
			return current, nil
		}
		return Ride{}, ErrInvalidTransition
	}

	trimmedDriverID := strings.TrimSpace(driverID)
	if trimmedDriverID == "" || trimmedDriverID == current.RiderID {
		return Ride{}, ErrUnauthorizedRideAction
	}

	updated, err := s.repository.UpdateStatusIfCurrent(ctx, id, StatusRequested, StatusAccepted, &trimmedDriverID)
	if err != nil {
		return Ride{}, err
	}

	return updated, s.publishRideEvent(ctx, "ride.accepted", updated)
}

func (s *Service) Start(ctx context.Context, id string, driverID string) (Ride, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return Ride{}, err
	}
	if current.Status != StatusAccepted {
		if current.Status == StatusStarted && driverOwnsRide(current, driverID) {
			return current, nil
		}
		return Ride{}, ErrInvalidTransition
	}
	if !driverOwnsRide(current, driverID) {
		return Ride{}, ErrUnauthorizedRideAction
	}

	updated, err := s.repository.UpdateStatusIfCurrent(ctx, id, StatusAccepted, StatusStarted, nil)
	if err != nil {
		return Ride{}, err
	}

	return updated, s.publishRideEvent(ctx, "ride.started", updated)
}

func (s *Service) Complete(ctx context.Context, id string, driverID string) (Ride, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return Ride{}, err
	}
	if current.Status != StatusStarted {
		if current.Status == StatusCompleted && driverOwnsRide(current, driverID) {
			return current, nil
		}
		return Ride{}, ErrInvalidTransition
	}
	if !driverOwnsRide(current, driverID) {
		return Ride{}, ErrUnauthorizedRideAction
	}

	updated, err := s.repository.UpdateStatusIfCurrent(ctx, id, StatusStarted, StatusCompleted, nil)
	if err != nil {
		return Ride{}, err
	}

	return updated, s.publishRideEvent(ctx, "ride.completed", updated)
}

func (s *Service) Cancel(ctx context.Context, id string, actorID string) (Ride, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return Ride{}, err
	}
	if current.Status == StatusCancelled && riderOrAssignedDriver(current, actorID) {
		return current, nil
	}
	if current.Status == StatusCompleted {
		return Ride{}, ErrInvalidTransition
	}
	if !riderOrAssignedDriver(current, actorID) {
		return Ride{}, ErrUnauthorizedRideAction
	}

	updated, err := s.repository.UpdateStatusIfCurrent(ctx, id, current.Status, StatusCancelled, nil)
	if err != nil {
		return Ride{}, err
	}

	return updated, s.publishRideEvent(ctx, "ride.cancelled", updated)
}

func (s *Service) publishRideEvent(ctx context.Context, eventType string, ride Ride) error {
	return s.publisher.Publish(ctx, events.Event{
		Type: eventType,
		Key:  ride.ID,
		Data: ride,
	})
}

func driverOwnsRide(ride Ride, driverID string) bool {
	return ride.DriverID != nil && strings.TrimSpace(driverID) == *ride.DriverID
}

func riderOrAssignedDriver(ride Ride, actorID string) bool {
	trimmedActorID := strings.TrimSpace(actorID)
	if trimmedActorID == "" {
		return false
	}
	if trimmedActorID == ride.RiderID {
		return true
	}

	return ride.DriverID != nil && trimmedActorID == *ride.DriverID
}

func validLocation(location Location) bool {
	return !math.IsNaN(location.Latitude) &&
		!math.IsNaN(location.Longitude) &&
		location.Latitude >= -90 &&
		location.Latitude <= 90 &&
		location.Longitude >= -180 &&
		location.Longitude <= 180
}
