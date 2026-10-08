package grpc

import (
	"context"
	"errors"
	"time"

	userv1 "github.com/ErenKarakus1/Ride-Sharing-System/backend/proto/gen/go/user/v1"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/user-service/internal/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	userv1.UnimplementedUserServiceServer
	service *user.Service
}

func NewUserServer(service *user.Service) *UserServer {
	return &UserServer{
		service: service,
	}
}

func (s *UserServer) CreateUser(ctx context.Context, request *userv1.CreateUserRequest) (*userv1.User, error) {
	created, err := s.service.CreateWithID(ctx, user.CreateUserInput{
		ID:          request.GetId(),
		Email:       request.GetEmail(),
		DisplayName: request.GetDisplayName(),
		PhoneNumber: request.GetPhoneNumber(),
		Role:        roleFromProto(request.GetRole()),
	})
	if errors.Is(err, user.ErrInvalidRole) {
		return nil, status.Error(codes.InvalidArgument, "invalid user role")
	}
	if errors.Is(err, user.ErrEmailAlreadyExists) {
		return nil, status.Error(codes.AlreadyExists, "email already exists")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return userToProto(created), nil
}

func (s *UserServer) GetUser(ctx context.Context, request *userv1.GetUserRequest) (*userv1.User, error) {
	found, err := s.service.Get(ctx, request.GetId())
	if errors.Is(err, user.ErrUserNotFound) {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return userToProto(found), nil
}

func roleFromProto(role userv1.Role) user.Role {
	switch role {
	case userv1.Role_ROLE_RIDER:
		return user.RoleRider
	case userv1.Role_ROLE_DRIVER:
		return user.RoleDriver
	default:
		return ""
	}
}

func roleToProto(role user.Role) userv1.Role {
	switch role {
	case user.RoleRider:
		return userv1.Role_ROLE_RIDER
	case user.RoleDriver:
		return userv1.Role_ROLE_DRIVER
	default:
		return userv1.Role_ROLE_UNSPECIFIED
	}
}

func userToProto(user user.User) *userv1.User {
	return &userv1.User{
		Id:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		PhoneNumber: user.PhoneNumber,
		Role:        roleToProto(user.Role),
		CreatedAt:   formatTime(user.CreatedAt),
		UpdatedAt:   formatTime(user.UpdatedAt),
	}
}

func formatTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339Nano)
}
