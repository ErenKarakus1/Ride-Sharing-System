package userclient

import (
	"context"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/auth-service/internal/auth"
	userv1 "github.com/ErenKarakus1/Ride-Sharing-System/backend/proto/gen/go/user/v1"
	"google.golang.org/grpc"
)

type GRPCProfileClient struct {
	client userv1.UserServiceClient
}

func NewGRPCProfileClient(conn *grpc.ClientConn) *GRPCProfileClient {
	return &GRPCProfileClient{
		client: userv1.NewUserServiceClient(conn),
	}
}

func (c *GRPCProfileClient) Create(ctx context.Context, input auth.CreateProfileInput) error {
	_, err := c.client.CreateUser(ctx, &userv1.CreateUserRequest{
		Id:          input.ID,
		Email:       input.Email,
		DisplayName: input.DisplayName,
		PhoneNumber: input.PhoneNumber,
		Role:        roleToProto(input.Role),
	})

	return err
}

func roleToProto(role string) userv1.Role {
	switch role {
	case "rider":
		return userv1.Role_ROLE_RIDER
	case "driver":
		return userv1.Role_ROLE_DRIVER
	default:
		return userv1.Role_ROLE_UNSPECIFIED
	}
}
