package payment

import (
	"context"
	"testing"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/payment-service/internal/events"
)

func TestPaymentTransitions(t *testing.T) {
	repository := newFakeRepository()
	service := NewService(repository, noopPublisher{})

	authorized, err := service.Authorize(context.Background(), AuthorizeRequest{
		RideID:  "ride-1",
		RiderID: "rider-1",
		Amount:  100,
	})
	if err != nil {
		t.Fatalf("authorize payment: %v", err)
	}
	if authorized.Status != StatusAuthorized {
		t.Fatalf("expected authorized, got %s", authorized.Status)
	}

	captured, err := service.Capture(context.Background(), authorized.ID)
	if err != nil {
		t.Fatalf("capture payment: %v", err)
	}
	if captured.Status != StatusCaptured {
		t.Fatalf("expected captured, got %s", captured.Status)
	}

	refunded, err := service.Refund(context.Background(), authorized.ID)
	if err != nil {
		t.Fatalf("refund payment: %v", err)
	}
	if refunded.Status != StatusRefunded {
		t.Fatalf("expected refunded, got %s", refunded.Status)
	}
}

func TestInvalidPaymentTransition(t *testing.T) {
	repository := newFakeRepository()
	service := NewService(repository, noopPublisher{})

	authorized, err := service.Authorize(context.Background(), AuthorizeRequest{
		RideID:  "ride-1",
		RiderID: "rider-1",
		Amount:  100,
	})
	if err != nil {
		t.Fatalf("authorize payment: %v", err)
	}

	if _, err := service.Refund(context.Background(), authorized.ID); err != ErrInvalidPaymentTransition {
		t.Fatalf("expected invalid transition, got %v", err)
	}
}

func TestAuthorizeRejectsInvalidPaymentInput(t *testing.T) {
	repository := newFakeRepository()
	service := NewService(repository, noopPublisher{})

	_, err := service.Authorize(context.Background(), AuthorizeRequest{
		RideID:  "ride-1",
		RiderID: "rider-1",
		Amount:  0,
	})
	if err != ErrInvalidPaymentAmount {
		t.Fatalf("expected invalid payment amount, got %v", err)
	}

	_, err = service.Authorize(context.Background(), AuthorizeRequest{
		RideID:   "ride-2",
		RiderID:  "rider-1",
		Amount:   10,
		Currency: "TURKISH_LIRA",
	})
	if err != ErrInvalidCurrency {
		t.Fatalf("expected invalid currency, got %v", err)
	}
}

func TestCompletedRideCapturesAuthorizedPayment(t *testing.T) {
	repository := newFakeRepository()
	service := NewService(repository, noopPublisher{})

	authorized, err := service.Authorize(context.Background(), AuthorizeRequest{
		RideID:  "ride-1",
		RiderID: "rider-1",
		Amount:  100,
	})
	if err != nil {
		t.Fatalf("authorize payment: %v", err)
	}

	err = service.HandleRideEvent(context.Background(), events.RideEvent{
		Type: "ride.completed",
		Data: events.RideEventData{ID: authorized.RideID},
	})
	if err != nil {
		t.Fatalf("handle ride event: %v", err)
	}

	captured, err := service.Get(context.Background(), authorized.ID)
	if err != nil {
		t.Fatalf("get payment: %v", err)
	}
	if captured.Status != StatusCaptured {
		t.Fatalf("expected captured, got %s", captured.Status)
	}
}

type fakeRepository struct {
	payments map[string]Payment
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{
		payments: make(map[string]Payment),
	}
}

func (r *fakeRepository) Create(ctx context.Context, payment Payment) (Payment, error) {
	payment.ID = "payment-1"
	r.payments[payment.ID] = payment
	return payment, nil
}

func (r *fakeRepository) Get(ctx context.Context, id string) (Payment, error) {
	payment, ok := r.payments[id]
	if !ok {
		return Payment{}, ErrPaymentNotFound
	}

	return payment, nil
}

func (r *fakeRepository) GetByRide(ctx context.Context, rideID string) (Payment, error) {
	for _, payment := range r.payments {
		if payment.RideID == rideID {
			return payment, nil
		}
	}

	return Payment{}, ErrPaymentNotFound
}

func (r *fakeRepository) GetAuthorizedByRide(ctx context.Context, rideID string) (Payment, error) {
	for _, payment := range r.payments {
		if payment.RideID == rideID && payment.Status == StatusAuthorized {
			return payment, nil
		}
	}

	return Payment{}, ErrPaymentNotFound
}

func (r *fakeRepository) UpdateStatus(ctx context.Context, id string, status Status) (Payment, error) {
	payment, ok := r.payments[id]
	if !ok {
		return Payment{}, ErrPaymentNotFound
	}

	payment.Status = status
	r.payments[id] = payment
	return payment, nil
}

type noopPublisher struct{}

func (noopPublisher) Publish(ctx context.Context, event events.Event) error {
	return nil
}
