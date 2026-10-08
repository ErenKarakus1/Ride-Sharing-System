package payment

import (
	"context"
	"errors"
)

var ErrPaymentNotFound = errors.New("payment not found")

type Repository interface {
	Create(ctx context.Context, payment Payment) (Payment, error)
	Get(ctx context.Context, id string) (Payment, error)
	UpdateStatus(ctx context.Context, id string, status Status) (Payment, error)
}
