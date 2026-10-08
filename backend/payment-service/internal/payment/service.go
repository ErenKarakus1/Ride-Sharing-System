package payment

import (
	"context"
	"errors"
	"strings"
)

var ErrInvalidPaymentTransition = errors.New("invalid payment status transition")

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Authorize(ctx context.Context, request AuthorizeRequest) (Payment, error) {
	currency := strings.TrimSpace(request.Currency)
	if currency == "" {
		currency = "TRY"
	}

	return s.repository.Create(ctx, Payment{
		RideID:   strings.TrimSpace(request.RideID),
		RiderID:  strings.TrimSpace(request.RiderID),
		DriverID: request.DriverID,
		Amount:   request.Amount,
		Currency: currency,
		Status:   StatusAuthorized,
	})
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

	return s.repository.UpdateStatus(ctx, id, StatusCaptured)
}

func (s *Service) Refund(ctx context.Context, id string) (Payment, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return Payment{}, err
	}
	if current.Status != StatusCaptured {
		return Payment{}, ErrInvalidPaymentTransition
	}

	return s.repository.UpdateStatus(ctx, id, StatusRefunded)
}
