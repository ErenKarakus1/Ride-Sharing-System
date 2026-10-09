package grpc

import (
	"context"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/auth-service/internal/auth"
	authv1 "github.com/ErenKarakus1/Ride-Sharing-System/backend/proto/gen/go/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	authv1.UnimplementedAuthServiceServer
	service *auth.Service
}

func NewAuthServer(service *auth.Service) *AuthServer {
	return &AuthServer{
		service: service,
	}
}

func (s *AuthServer) ValidateToken(ctx context.Context, request *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	claims, err := s.service.ValidateToken(request.GetAccessToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}

	return &authv1.ValidateTokenResponse{
		Valid:  true,
		UserId: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}, nil
}
