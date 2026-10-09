package config

import "testing"

func TestValidateRejectsUnsafeJWTSecretOutsideLocal(t *testing.T) {
	cfg := Config{
		Environment: "production",
		JWTSecret:   "change-me",
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected invalid production JWT secret")
	}
}

func TestValidateAllowsLocalDefaultSecret(t *testing.T) {
	cfg := Config{
		Environment: "local",
		JWTSecret:   "change-me",
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected local default secret to be accepted, got %v", err)
	}
}
