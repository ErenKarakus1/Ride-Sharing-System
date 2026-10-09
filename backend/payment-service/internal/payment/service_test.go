package payment

import (
	"context"
	"testing"
)

func TestPaymentTransitions(t *testing.T) {
	repository := newFakeRepository()
	service := NewService(repository)

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
	service := NewService(repository)

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

func (r *fakeRepository) UpdateStatus(ctx context.Context, id string, status Status) (Payment, error) {
	payment, ok := r.payments[id]
	if !ok {
		return Payment{}, ErrPaymentNotFound
	}

	payment.Status = status
	r.payments[id] = payment
	return payment, nil
}
