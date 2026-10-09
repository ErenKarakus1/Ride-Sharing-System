package payment

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/payment-service/internal/events"
)

var ErrInvalidPaymentTransition = errors.New("invalid payment status transition")

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

func (s *Service) Authorize(ctx context.Context, request AuthorizeRequest) (Payment, error) {
	existing, err := s.repository.GetByRide(ctx, request.RideID)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, ErrPaymentNotFound) {
		return Payment{}, err
	}

	currency := strings.TrimSpace(request.Currency)
	if currency == "" {
		currency = "TRY"
	}

	payment, err := s.repository.Create(ctx, Payment{
		RideID:   strings.TrimSpace(request.RideID),
		RiderID:  strings.TrimSpace(request.RiderID),
		DriverID: request.DriverID,
		Amount:   request.Amount,
		Currency: currency,
		Status:   StatusAuthorized,
	})
	if err != nil {
		return Payment{}, err
	}

	return payment, s.publishPaymentEvent(ctx, "payment.authorized", payment)
}

func (s *Service) Get(ctx context.Context, id string) (Payment, error) {
	return s.repository.Get(ctx, strings.TrimSpace(id))
}

func (s *Service) Capture(ctx context.Context, id string) (Payment, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return Payment{}, err
	}
	if current.Status != StatusAuthorized {
		return Payment{}, ErrInvalidPaymentTransition
	}

	payment, err := s.repository.UpdateStatus(ctx, id, StatusCaptured)
	if err != nil {
		return Payment{}, err
	}

	return payment, s.publishPaymentEvent(ctx, "payment.captured", payment)
}

func (s *Service) Refund(ctx context.Context, id string) (Payment, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return Payment{}, err
	}
	if current.Status != StatusCaptured {
		return Payment{}, ErrInvalidPaymentTransition
	}

	payment, err := s.repository.UpdateStatus(ctx, id, StatusRefunded)
	if err != nil {
		return Payment{}, err
	}

	return payment, s.publishPaymentEvent(ctx, "payment.refunded", payment)
}

func (s *Service) HandleRideEvent(ctx context.Context, event events.RideEvent) error {
	if event.Type != "ride.completed" {
		return nil
	}

	payment, err := s.repository.GetAuthorizedByRide(ctx, event.Data.ID)
	if errors.Is(err, ErrPaymentNotFound) {
		log.Printf("no authorized payment found for completed ride_id=%s", event.Data.ID)
		return nil
	}
	if err != nil {
		return err
	}

	_, err = s.Capture(ctx, payment.ID)
	return err
}

func (s *Service) publishPaymentEvent(ctx context.Context, eventType string, payment Payment) error {
	if s.publisher == nil {
		return nil
	}

	return s.publisher.Publish(ctx, events.Event{
		Type: eventType,
		Key:  payment.ID,
		Data: payment,
	})
}
