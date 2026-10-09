package auth

import (
	"context"
	"testing"
)

func TestRegisterRejectsWeakPassword(t *testing.T) {
	service := NewService(noopRepository{}, NewTokenIssuer("secret", 0), noopProfiles{})

	_, err := service.Register(context.Background(), RegisterRequest{
		Email:       "rider@example.com",
		Password:    "password123",
		DisplayName: "Rider",
		PhoneNumber: "+905551112233",
		Role:        "rider",
	})
	if err != ErrWeakPassword {
		t.Fatalf("expected weak password, got %v", err)
	}
}

func TestStrongPassword(t *testing.T) {
	if !strongPassword("Password12345") {
		t.Fatal("expected password to be strong")
	}
	if strongPassword("password12345") {
		t.Fatal("expected missing uppercase password to be weak")
	}
	if strongPassword("PASSWORD12345") {
		t.Fatal("expected missing lowercase password to be weak")
	}
	if strongPassword("PasswordOnly") {
		t.Fatal("expected missing digit password to be weak")
	}
}

type noopRepository struct{}

func (noopRepository) Create(ctx context.Context, account Account) (Account, error) {
	return account, nil
}

func (noopRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (noopRepository) GetByEmail(ctx context.Context, email string) (Account, error) {
	return Account{}, ErrAccountNotFound
}

type noopProfiles struct{}

func (noopProfiles) Create(ctx context.Context, input CreateProfileInput) error {
	return nil
}
