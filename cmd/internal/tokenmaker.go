package internal

import (
	"fmt"
	"log"
	"rbac/internal/envvar"
	"rbac/internal/tokenmaker"
	"rbac/internal/tokenmaker/jwtmaker"
	"rbac/internal/tokenmaker/pasetomaker"
	"strconv"
	"time"
)

func NewTokenMaker(conf *envvar.Configuration) (tokenmaker.TokenMaker, error) {
	get := func(v string) string {
		res, err := conf.Get(v)
		if err != nil {
			log.Fatalf("Couldn't get configuration value for %s: %s", v, err)
		}

		return res
	}
	tokenSymmetricKey := get("TOKEN_SYMMETRIC_KEY")
	tokenExpiration := get("TOKEN_EXPIRATION")
	tokenMaker := get("TOKEN_MAKER")
	var token_maker tokenmaker.TokenMaker
	duration, err := strconv.Atoi(tokenExpiration)
	if err != nil {
		return nil, fmt.Errorf("invalid durration: %s", err)
	}
	maker, err := pasetomaker.NewPasetoMaker(tokenSymmetricKey, time.Duration(duration)*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("couldn't create token maker: %s", err)
	}
	token_maker = maker

	if tokenMaker == "JWT" {
		maker, err := jwtmaker.NewJWTMaker(tokenSymmetricKey, time.Duration(duration)*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("couldn't create token maker: %s", err)
		}
		token_maker = maker
	}
	return token_maker, nil

}
