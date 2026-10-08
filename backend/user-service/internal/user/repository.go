package user

import (
	"context"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

type Repository interface {
	Create(ctx context.Context, user User) (User, error)
	List(ctx context.Context) ([]User, error)
	Get(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, id string, request UpdateUserRequest) (User, error)
}
