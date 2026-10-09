package ride

import (
	"context"
	"testing"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/ride-service/internal/events"
)

func TestRideTransitions(t *testing.T) {
	repository := newFakeRepository()
	service := NewService(repository, noopPublisher{})

	created, err := service.Create(context.Background(), CreateRideRequest{
		RiderID: "rider-1",
		Pickup:  Location{Latitude: 1, Longitude: 2, Address: "pickup"},
		Dropoff: Location{Latitude: 3, Longitude: 4, Address: "dropoff"},
	})
	if err != nil {
		t.Fatalf("create ride: %v", err)
	}

	accepted, err := service.Accept(context.Background(), created.ID, "driver-1")
	if err != nil {
		t.Fatalf("accept ride: %v", err)
	}
	if accepted.Status != StatusAccepted {
		t.Fatalf("expected accepted, got %s", accepted.Status)
	}

	started, err := service.Start(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("start ride: %v", err)
	}
	if started.Status != StatusStarted {
		t.Fatalf("expected started, got %s", started.Status)
	}

	completed, err := service.Complete(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("complete ride: %v", err)
	}
	if completed.Status != StatusCompleted {
		t.Fatalf("expected completed, got %s", completed.Status)
	}
}

func TestInvalidRideTransition(t *testing.T) {
	repository := newFakeRepository()
	service := NewService(repository, noopPublisher{})

	created, err := service.Create(context.Background(), CreateRideRequest{
		RiderID: "rider-1",
		Pickup:  Location{Latitude: 1, Longitude: 2, Address: "pickup"},
		Dropoff: Location{Latitude: 3, Longitude: 4, Address: "dropoff"},
	})
	if err != nil {
		t.Fatalf("create ride: %v", err)
	}

	if _, err := service.Start(context.Background(), created.ID); err != ErrInvalidTransition {
		t.Fatalf("expected invalid transition, got %v", err)
	}
}

type noopPublisher struct{}

func (noopPublisher) Publish(ctx context.Context, event events.Event) error {
	return nil
}

type fakeRepository struct {
	rides map[string]Ride
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{
		rides: make(map[string]Ride),
	}
}

func (r *fakeRepository) Create(ctx context.Context, ride Ride) (Ride, error) {
	ride.ID = "ride-1"
	ride.Status = StatusRequested
	r.rides[ride.ID] = ride
	return ride, nil
}

func (r *fakeRepository) Get(ctx context.Context, id string) (Ride, error) {
	ride, ok := r.rides[id]
	if !ok {
		return Ride{}, ErrRideNotFound
	}

	return ride, nil
}

func (r *fakeRepository) ListByRider(ctx context.Context, riderID string) ([]Ride, error) {
	return nil, nil
}

func (r *fakeRepository) UpdateStatus(ctx context.Context, id string, status Status, driverID *string) (Ride, error) {
	ride, ok := r.rides[id]
	if !ok {
		return Ride{}, ErrRideNotFound
	}

	ride.Status = status
	if driverID != nil {
		ride.DriverID = driverID
	}
	r.rides[id] = ride
	return ride, nil
}
