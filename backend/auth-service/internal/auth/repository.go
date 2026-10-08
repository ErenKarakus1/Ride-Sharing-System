package auth

import (
	"context"
	"errors"
)

var ErrAccountNotFound = errors.New("account not found")
var ErrEmailAlreadyExists = errors.New("email already exists")

type Repository interface {
	Create(ctx context.Context, account Account) (Account, error)
	Delete(ctx context.Context, id string) error
	GetByEmail(ctx context.Context, email string) (Account, error)
}
