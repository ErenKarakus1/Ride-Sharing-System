package auth

import (
	"context"
	"errors"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidRole = errors.New("invalid role")
var ErrWeakPassword = errors.New("weak password")

type Service struct {
	repository  Repository
	tokenIssuer *TokenIssuer
	profiles    ProfileClient
}

func NewService(repository Repository, tokenIssuer *TokenIssuer, profiles ProfileClient) *Service {
	return &Service{
		repository:  repository,
		tokenIssuer: tokenIssuer,
		profiles:    profiles,
	}
}

func (s *Service) Register(ctx context.Context, request RegisterRequest) (AuthResponse, error) {
	if !strongPassword(request.Password) {
		return AuthResponse{}, ErrWeakPassword
	}

	role := strings.TrimSpace(request.Role)
	if role != "rider" && role != "driver" {
		return AuthResponse{}, ErrInvalidRole
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResponse{}, err
	}

	account, err := s.repository.Create(ctx, Account{
		Email:        strings.TrimSpace(request.Email),
		PasswordHash: string(passwordHash),
		Role:         role,
	})
	if err != nil {
		return AuthResponse{}, err
	}

	if err := s.profiles.Create(ctx, CreateProfileInput{
		ID:          account.ID,
		Email:       account.Email,
		DisplayName: request.DisplayName,
		PhoneNumber: request.PhoneNumber,
		Role:        request.Role,
	}); err != nil {
		_ = s.repository.Delete(ctx, account.ID)
		return AuthResponse{}, err
	}

	return s.authResponse(account)
}

func strongPassword(password string) bool {
	if len(password) < 12 {
		return false
	}

	var hasLower bool
	var hasUpper bool
	var hasDigit bool
	for _, value := range password {
		hasLower = hasLower || unicode.IsLower(value)
		hasUpper = hasUpper || unicode.IsUpper(value)
		hasDigit = hasDigit || unicode.IsDigit(value)
	}

	return hasLower && hasUpper && hasDigit
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

func (s *Service) ValidateToken(accessToken string) (TokenClaims, error) {
	return s.tokenIssuer.Validate(accessToken)
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
		Role:        account.Role,
	}, nil
}
