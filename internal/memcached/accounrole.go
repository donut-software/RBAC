package memcached

import (
	"bytes"
	"context"
	"encoding/gob"
	"rbac/internal"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"go.uber.org/zap"
)

// Index ...
func (t *RBAC) IndexAccountRole(ctx context.Context, accRole internal.AccountRoles) error {
	return t.orig.IndexAccountRole(ctx, accRole)
}

func (t *RBAC) GetAccountRole(ctx context.Context, accRoleId string) (internal.AccountRoles, error) {
	key := "accountrole_" + accRoleId
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			res, err := t.orig.GetAccountRole(ctx, accRoleId)
			if err != nil {
				return internal.AccountRoles{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetAccountRole")
			}
			account, err := t.orig.GetAccount(ctx, res.Account.UserName)
			if err != nil {
				return internal.AccountRoles{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetAccount")
			}
			account.Profile, err = t.orig.GetProfile(ctx, account.Profile.Id)
			if err != nil {
				return internal.AccountRoles{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetProfile")
			}
			res.Account = account
			role, err := t.orig.GetRole(ctx, res.Role.Id)
			if err != nil {
				return internal.AccountRoles{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetRole")
			}
			res.Role = role
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
		return internal.AccountRoles{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.AccountRoles
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.AccountRoles{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}
func (t *RBAC) DeleteAccountRole(ctx context.Context, accRoleId string) error {
	return t.orig.DeleteAccountRole(ctx, accRoleId)
}

func (t *RBAC) UpdateAccountRole(ctx context.Context, accountRole internal.AccountRoles) error {
	err := t.orig.DeleteAccountRole(ctx, accountRole.Id)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.DeleteAccountRole")
	}
	return t.orig.IndexAccountRole(ctx, accountRole)
}

func (t *RBAC) GetAccountRoleByAccount(ctx context.Context, username string) (internal.AccountRoleByAccountResult, error) {
	key := "accountrolebyaccount_" + username
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			res, err := t.orig.AccountRoleByAccount(ctx, &username)
			if err != nil {
				return internal.AccountRoleByAccountResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.AccountRoleByAccount")
			}
			acc, err := t.orig.GetAccount(ctx, res.Account.UserName)
			if err != nil {
				return internal.AccountRoleByAccountResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetAccount")
			}
			acc.Profile, err = t.orig.GetProfile(ctx, acc.Profile.Id)
			if err != nil {
				return internal.AccountRoleByAccountResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetProfile")
			}
			res.Account = acc
			roles := []internal.Roles{}
			for _, value := range res.Roles {
				role, err := t.orig.GetRole(ctx, value.Id)
				if err != nil {
					return internal.AccountRoleByAccountResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetRole")
				}
				roles = append(roles, role)
			}
			res.Roles = roles
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
		return internal.AccountRoleByAccountResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.AccountRoleByAccountResult
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.AccountRoleByAccountResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}

func (t *RBAC) GetAccountRoleByRole(ctx context.Context, roleid string) (internal.AccountRoleByRoleResult, error) {
	key := "accountrolebyrole_" + roleid
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			res, err := t.orig.AccountRoleByRole(ctx, &roleid)
			if err != nil {
				return internal.AccountRoleByRoleResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.AccountRoleByRole")
			}
			role, err := t.orig.GetRole(ctx, res.Role.Id)
			if err != nil {
				return internal.AccountRoleByRoleResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetRole")
			}
			res.Role = role
			accounts := []internal.Account{}
			for _, value := range res.Account {
				acc, err := t.orig.GetAccount(ctx, value.UserName)
				if err != nil {
					return internal.AccountRoleByRoleResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetAccount")
				}
				acc.Profile, err = t.orig.GetProfile(ctx, acc.Profile.Id)
				if err != nil {
					return internal.AccountRoleByRoleResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetProfile")
				}
				accounts = append(accounts, acc)
			}
			res.Account = accounts
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
		return internal.AccountRoleByRoleResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.AccountRoleByRoleResult
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.AccountRoleByRoleResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}

func (t *RBAC) ListAccountRole(ctx context.Context, args internal.ListArgs) (internal.ListAccountRole, error) {
	key := newKey("listaccountrole", args)
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			listacc, err := t.orig.ListAccountRole(ctx, args)
			if err != nil {
				return internal.ListAccountRole{}, err
			}
			for i, value := range listacc.AccountRoles {
				acc, err := t.GetAccount(ctx, value.Account.UserName)
				if err != nil {
					return internal.ListAccountRole{}, err
				}
				rl, err := t.GetRole(ctx, value.Role.Id)
				if err != nil {
					return internal.ListAccountRole{}, err
				}
				listacc.AccountRoles[i].Account = acc
				listacc.AccountRoles[i].Role = rl
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
		return internal.ListAccountRole{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.ListAccountRole
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.ListAccountRole{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}
