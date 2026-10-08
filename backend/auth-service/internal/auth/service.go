package auth

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service struct {
	repository  Repository
	tokenIssuer *TokenIssuer
}

func NewService(repository Repository, tokenIssuer *TokenIssuer) *Service {
	return &Service{
		repository:  repository,
		tokenIssuer: tokenIssuer,
	}
}

func (s *Service) Register(ctx context.Context, request RegisterRequest) (AuthResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResponse{}, err
	}

	account, err := s.repository.Create(ctx, Account{
		Email:        strings.TrimSpace(request.Email),
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		return AuthResponse{}, err
	}

	return s.authResponse(account)
}

func (s *Service) Login(ctx context.Context, request LoginRequest) (AuthResponse, error) {
	account, err := s.repository.GetByEmail(ctx, request.Email)
	if errors.Is(err, ErrAccountNotFound) {
		return AuthResponse{}, ErrInvalidCredentials
	}
	if err != nil {
		return AuthResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(request.Password)); err != nil {
		return AuthResponse{}, ErrInvalidCredentials
	}

	return s.authResponse(account)
}

func (s *Service) authResponse(account Account) (AuthResponse, error) {
	accessToken, expiresAt, err := s.tokenIssuer.Issue(account)
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt,
		UserID:      account.ID,
		Email:       account.Email,
	}, nil
}
