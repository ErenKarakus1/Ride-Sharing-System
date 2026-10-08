package auth

import (
	"errors"
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

func (i *TokenIssuer) Validate(accessToken string) (TokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}

		return i.secret, nil
	})
	if err != nil {
		return TokenClaims{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return TokenClaims{}, errors.New("invalid token")
	}

	userID, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)
	if userID == "" || email == "" {
		return TokenClaims{}, errors.New("missing token claims")
	}

	return TokenClaims{
		UserID: userID,
		Email:  email,
	}, nil
}

type TokenClaims struct {
	UserID string
	Email  string
}
