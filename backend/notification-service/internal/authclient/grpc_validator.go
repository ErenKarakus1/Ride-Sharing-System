package authclient

import (
	"context"

	authv1 "github.com/ErenKarakus1/Ride-Sharing-System/backend/proto/gen/go/auth/v1"
	"google.golang.org/grpc"
)

type Claims struct {
	UserID string
	Email  string
}

type Validator interface {
	Validate(ctx context.Context, accessToken string) (Claims, error)
}

type GRPCValidator struct {
	client authv1.AuthServiceClient
}

func NewGRPCValidator(conn *grpc.ClientConn) *GRPCValidator {
	return &GRPCValidator{
		client: authv1.NewAuthServiceClient(conn),
	}
}

func (v *GRPCValidator) Validate(ctx context.Context, accessToken string) (Claims, error) {
	response, err := v.client.ValidateToken(ctx, &authv1.ValidateTokenRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		return Claims{}, err
	}

	return Claims{
		UserID: response.GetUserId(),
		Email:  response.GetEmail(),
	}, nil
}
