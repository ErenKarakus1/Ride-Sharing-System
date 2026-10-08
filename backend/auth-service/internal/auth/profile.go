package auth

import "context"

type CreateProfileInput struct {
	ID          string
	Email       string
	DisplayName string
	PhoneNumber string
	Role        string
}

type ProfileClient interface {
	Create(ctx context.Context, input CreateProfileInput) error
}
