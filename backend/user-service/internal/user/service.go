package user

import (
	"context"
	"errors"
	"strings"
)

var ErrInvalidRole = errors.New("invalid user role")

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(ctx context.Context, request CreateUserRequest) (User, error) {
	user := request.ToUser()
	if !user.Role.Valid() {
		return User{}, ErrInvalidRole
	}

	return s.repository.Create(ctx, user)
}

func (s *Service) List(ctx context.Context) ([]User, error) {
	return s.repository.List(ctx)
}

func (s *Service) Get(ctx context.Context, id string) (User, error) {
	return s.repository.Get(ctx, strings.TrimSpace(id))
}

func (s *Service) Update(ctx context.Context, id string, request UpdateUserRequest) (User, error) {
	return s.repository.Update(ctx, strings.TrimSpace(id), request.Trimmed())
}
