package memcached

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"rbac/internal"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"go.uber.org/zap"
)

// Index ...
func (t *RBAC) IndexAccount(ctx context.Context, account internal.Account) error {
	err := t.orig.IndexProfile(ctx, account.Profile)
	if err != nil {
		return err
	}
	return t.orig.IndexAccount(ctx, account)
}

func (t *RBAC) GetAccount(ctx context.Context, username string) (internal.Account, error) {
	key := "account_" + username
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			res, err := t.orig.GetAccount(ctx, username)
			if err != nil {
				return internal.Account{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetAccount")
			}
			prof, err := t.orig.GetProfile(ctx, res.Profile.Id)
			if err != nil {
				return internal.Account{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetProfile")
			}
			res.Profile = prof
			var b bytes.Buffer
			if err := gob.NewEncoder(&b).Encode(&res); err == nil {
				t.logger.Info("settin value")

				t.client.Set(&memcache.Item{
					Key:        key,
					Value:      b.Bytes(),
					Expiration: int32(time.Now().Add(25 * time.Second).Unix()),
				})
			}

			return res, err
		}
		return internal.Account{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.Account
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.Account{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}
func (t *RBAC) DeleteAccount(ctx context.Context, username string, profileId string) error {
	err := t.orig.DeleteProfile(ctx, profileId)
	if err != nil {
		return err
	}
	return t.orig.DeleteAccount(ctx, username)
}
func (t *RBAC) ListAccount(ctx context.Context, args internal.ListAccountArgs) (internal.ListAccount, error) {
	key := newKey(args)
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			listacc, err := t.orig.ListAccount(ctx, args)
			if err != nil {
				return internal.ListAccount{}, err
			}
			for i, value := range listacc.Accounts {
				listacc.Accounts[i].Profile, err = t.orig.GetProfile(ctx, value.Profile.Id)
				if err != nil {
					return internal.ListAccount{}, err
				}
			}
			var b bytes.Buffer
			if err := gob.NewEncoder(&b).Encode(&listacc); err == nil {
				t.logger.Info("settin value")

				t.client.Set(&memcache.Item{
					Key:        key,
					Value:      b.Bytes(),
					Expiration: int32(time.Now().Add(25 * time.Second).Unix()),
				})
			}

			return listacc, err
		}
		return internal.ListAccount{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.ListAccount
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.ListAccount{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}

func newKey(args internal.ListAccountArgs) string {
	var (
		from int
		size int
	)

	if args.From != nil {
		from = *args.From
	}
	if args.Size != nil {
		size = *args.Size
	}

	// if args.Role != nil {
	// 	role = *args.Role
	// }

	return fmt.Sprintf("listaccount_%d_%d", from, size)
}
