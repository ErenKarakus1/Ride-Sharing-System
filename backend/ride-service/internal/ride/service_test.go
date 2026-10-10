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

	started, err := service.Start(context.Background(), created.ID, "driver-1")
	if err != nil {
		t.Fatalf("start ride: %v", err)
	}
	if started.Status != StatusStarted {
		t.Fatalf("expected started, got %s", started.Status)
	}

	completed, err := service.Complete(context.Background(), created.ID, "driver-1")
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

	if _, err := service.Start(context.Background(), created.ID, "driver-1"); err != ErrInvalidTransition {
		t.Fatalf("expected invalid transition, got %v", err)
	}
}

func TestCreateRejectsInvalidLocation(t *testing.T) {
	repository := newFakeRepository()
	service := NewService(repository, noopPublisher{})

	_, err := service.Create(context.Background(), CreateRideRequest{
		RiderID: "rider-1",
		Pickup:  Location{Latitude: 100, Longitude: 29},
		Dropoff: Location{Latitude: 41.1, Longitude: 29.1},
	})
	if err != ErrInvalidLocation {
		t.Fatalf("expected invalid location, got %v", err)
	}
}

func TestOnlyAssignedDriverCanAdvanceRide(t *testing.T) {
	repository := newFakeRepository()
	service := NewService(repository, noopPublisher{})

	created, err := service.Create(context.Background(), CreateRideRequest{
		RiderID: "rider-1",
		Pickup:  Location{Latitude: 41.0, Longitude: 29.0},
		Dropoff: Location{Latitude: 41.1, Longitude: 29.1},
	})
	if err != nil {
		t.Fatalf("create ride: %v", err)
	}

	if _, err := service.Accept(context.Background(), created.ID, "rider-1"); err != ErrUnauthorizedRideAction {
		t.Fatalf("expected unauthorized self-accept, got %v", err)
	}

	accepted, err := service.Accept(context.Background(), created.ID, "driver-1")
	if err != nil {
		t.Fatalf("accept ride: %v", err)
	}
	if _, err := service.Start(context.Background(), accepted.ID, "driver-2"); err != ErrUnauthorizedRideAction {
		t.Fatalf("expected unauthorized driver start, got %v", err)
	}
}

func TestSecondDriverCannotAcceptAlreadyAcceptedRide(t *testing.T) {
	repository := newFakeRepository()
	service := NewService(repository, noopPublisher{})

	created, err := service.Create(context.Background(), CreateRideRequest{
		RiderID: "rider-1",
		Pickup:  Location{Latitude: 41.0, Longitude: 29.0},
		Dropoff: Location{Latitude: 41.1, Longitude: 29.1},
	})
	if err != nil {
		t.Fatalf("create ride: %v", err)
	}

	if _, err := service.Accept(context.Background(), created.ID, "driver-1"); err != nil {
		t.Fatalf("accept ride: %v", err)
	}
	if _, err := service.Accept(context.Background(), created.ID, "driver-2"); err != ErrInvalidTransition {
		t.Fatalf("expected invalid transition for second accept, got %v", err)
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
	ride.ID = "11111111-1111-4111-8111-111111111111"
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

func (r *fakeRepository) UpdateStatusIfCurrent(ctx context.Context, id string, currentStatus Status, nextStatus Status, driverID *string) (Ride, error) {
	ride, ok := r.rides[id]
	if !ok {
		return Ride{}, ErrRideNotFound
	}
	if ride.Status != currentStatus {
		return Ride{}, ErrInvalidTransition
	}

	ride.Status = nextStatus
	if driverID != nil {
		ride.DriverID = driverID
	}
	r.rides[id] = ride
	return ride, nil
}
