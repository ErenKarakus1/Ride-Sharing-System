package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenIssuer struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenIssuer(secret string, ttl time.Duration) *TokenIssuer {
	return &TokenIssuer{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (i *TokenIssuer) Issue(account Account) (string, time.Time, error) {
	expiresAt := time.Now().UTC().Add(i.ttl)

	claims := jwt.MapClaims{
		"sub":   account.ID,
		"email": account.Email,
		"exp":   expiresAt.Unix(),
		"iat":   time.Now().UTC().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(i.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}
