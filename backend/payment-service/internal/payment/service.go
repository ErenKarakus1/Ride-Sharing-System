package payment

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/payment-service/internal/events"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/payment-service/internal/rideclient"
)

var ErrInvalidPaymentTransition = errors.New("invalid payment status transition")
var ErrInvalidPaymentAmount = errors.New("invalid payment amount")
var ErrInvalidCurrency = errors.New("invalid currency")
var ErrRidePaymentNotAllowed = errors.New("ride is not payable")
var ErrRideRiderMismatch = errors.New("payment rider does not match ride rider")

type Service struct {
	repository Repository
	publisher  events.Publisher
	rides      rideclient.Client
}

func NewService(repository Repository, publisher events.Publisher, rides rideclient.Client) *Service {
	return &Service{
		repository: repository,
		publisher:  publisher,
		rides:      rides,
	}
}

func (s *Service) Authorize(ctx context.Context, request AuthorizeRequest) (Payment, error) {
	if err := s.verifyRideForPayment(ctx, request); err != nil {
		return Payment{}, err
	}

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
	currency = strings.ToUpper(currency)
	if request.Amount <= 0 {
		return Payment{}, ErrInvalidPaymentAmount
	}
	if len(currency) != 3 {
		return Payment{}, ErrInvalidCurrency
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

func (s *Service) verifyRideForPayment(ctx context.Context, request AuthorizeRequest) error {
	if s.rides == nil {
		return nil
	}

	ride, err := s.rides.GetRide(ctx, strings.TrimSpace(request.RideID))
	if err != nil {
		return err
	}
	if ride.RiderID != strings.TrimSpace(request.RiderID) {
		return ErrRideRiderMismatch
	}
	if ride.Status == "completed" || ride.Status == "cancelled" {
		return ErrRidePaymentNotAllowed
	}

	return nil
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
