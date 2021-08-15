package tokenmaker

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TokenMaker interface {
	CreateToken(username string) (string, error)
	VerifyToken(token string) (*Payload, error)
}
type Account struct {
	tokenMaker TokenMaker
}

func NewTokenMaker(tokenMaker TokenMaker) *Account {
	return &Account{
		tokenMaker: tokenMaker,
	}
}

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
