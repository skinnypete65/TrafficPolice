package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenManager interface {
	NewJWT(userId string, ttl time.Duration) (string, error)
	Parse(accessToken string) (string, error)
}

type manager struct {
	signingKey string
}

func NewTokenManager(signingKey string) (TokenManager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &manager{signingKey: signingKey}, nil
}

func (m *manager) NewJWT(userId string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(ttl).Unix(),
		"sub": userId,
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *manager) Parse(accessToken string) (string, error) {
	keyFunc := func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	}
	token, err := jwt.Parse(accessToken, keyFunc)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["sub"].(string), nil
}
