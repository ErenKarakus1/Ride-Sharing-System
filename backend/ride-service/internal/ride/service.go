package ride

import (
	"context"
	"errors"
	"strings"
)

var ErrInvalidTransition = errors.New("invalid ride status transition")

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, request CreateRideRequest) (Ride, error) {
	return s.repository.Create(ctx, Ride{
		RiderID: strings.TrimSpace(request.RiderID),
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
		return Ride{}, ErrInvalidTransition
	}

	trimmedDriverID := strings.TrimSpace(driverID)
	return s.repository.UpdateStatus(ctx, id, StatusAccepted, &trimmedDriverID)
}

func (s *Service) Start(ctx context.Context, id string) (Ride, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return Ride{}, err
	}
	if current.Status != StatusAccepted {
		return Ride{}, ErrInvalidTransition
	}

	return s.repository.UpdateStatus(ctx, id, StatusStarted, nil)
}

func (s *Service) Complete(ctx context.Context, id string) (Ride, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return Ride{}, err
	}
	if current.Status != StatusStarted {
		return Ride{}, ErrInvalidTransition
	}

	return s.repository.UpdateStatus(ctx, id, StatusCompleted, nil)
}

func (s *Service) Cancel(ctx context.Context, id string) (Ride, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return Ride{}, err
	}
	if current.Status == StatusCompleted || current.Status == StatusCancelled {
		return Ride{}, ErrInvalidTransition
	}

	return s.repository.UpdateStatus(ctx, id, StatusCancelled, nil)
}
