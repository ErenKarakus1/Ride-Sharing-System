package auth

import "testing"
import "time"

func TestTokenIssuerValidateIncludesRole(t *testing.T) {
	issuer := NewTokenIssuer("secret", time.Hour)

	token, _, err := issuer.Issue(Account{
		ID:    "user-1",
		Email: "rider@example.com",
		Role:  "rider",
	})
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}

	claims, err := issuer.Validate(token)
	if err != nil {
		t.Fatalf("validate token: %v", err)
	}

	if claims.UserID != "user-1" || claims.Email != "rider@example.com" || claims.Role != "rider" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}
