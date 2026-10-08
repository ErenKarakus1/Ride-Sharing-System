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
	return s.CreateWithID(ctx, request.ToInput())
}

func (s *Service) CreateWithID(ctx context.Context, input CreateUserInput) (User, error) {
	if !input.Role.Valid() {
		return User{}, ErrInvalidRole
	}

	return s.repository.Create(ctx, User{
		ID:          strings.TrimSpace(input.ID),
		Email:       strings.TrimSpace(input.Email),
		DisplayName: strings.TrimSpace(input.DisplayName),
		PhoneNumber: strings.TrimSpace(input.PhoneNumber),
		Role:        input.Role,
	})
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
