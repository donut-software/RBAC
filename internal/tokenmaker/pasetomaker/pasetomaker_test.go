package pasetomaker_test

import (
	"rbac/internal/tokenmaker"
	"rbac/internal/tokenmaker/pasetomaker"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	td, err := strconv.Atoi("1")
	maker, err := pasetomaker.NewPasetoMaker(tokenmaker.RandomString(32), time.Duration(td)*time.Minute)
	require.NoError(t, err)

	username := tokenmaker.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := pasetomaker.NewPasetoMaker(tokenmaker.RandomString(32), -time.Minute)
	require.NoError(t, err)
	token, err := maker.CreateToken(tokenmaker.RandomOwner())
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, tokenmaker.ErrExpiredToken.Error())
	require.Nil(t, payload)
}
