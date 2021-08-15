package pasetomaker

import (
	"fmt"
	"rbac/internal/tokenmaker"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
	duration     time.Duration
}

func NewPasetoMaker(symmetricKey string, duration time.Duration) (tokenmaker.TokenMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
		duration:     duration,
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string) (string, error) {
	payload, err := tokenmaker.NewPayload(username, maker.duration)
	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}
func (maker *PasetoMaker) VerifyToken(token string) (*tokenmaker.Payload, error) {
	payload := &tokenmaker.Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, tokenmaker.ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}
