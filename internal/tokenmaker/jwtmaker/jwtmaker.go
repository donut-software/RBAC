package jwtmaker

import (
	"errors"
	"fmt"
	"rbac/internal/tokenmaker"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
	duration  time.Duration
}

func NewJWTMaker(secretKey string, duration time.Duration) (tokenmaker.TokenMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at lease %d characters", minSecretKeySize)
	}
	return &JWTMaker{
		secretKey: secretKey,
		duration:  duration,
	}, nil
}

func (maker *JWTMaker) CreateToken(username string) (string, error) {
	payload, err := tokenmaker.NewPayload(username, maker.duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}
func (maker *JWTMaker) VerifyToken(token string) (*tokenmaker.Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, tokenmaker.ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &tokenmaker.Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, tokenmaker.ErrExpiredToken) {
			return nil, tokenmaker.ErrExpiredToken
		}
		return nil, tokenmaker.ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*tokenmaker.Payload)
	if !ok {
		return nil, tokenmaker.ErrInvalidToken
	}
	return payload, nil
}
