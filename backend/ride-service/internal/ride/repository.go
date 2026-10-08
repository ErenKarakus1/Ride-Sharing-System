package ride

import (
	"context"
	"errors"
)

var ErrRideNotFound = errors.New("ride not found")

type Repository interface {
	Create(ctx context.Context, ride Ride) (Ride, error)
	Get(ctx context.Context, id string) (Ride, error)
	ListByRider(ctx context.Context, riderID string) ([]Ride, error)
	UpdateStatus(ctx context.Context, id string, status Status, driverID *string) (Ride, error)
}
