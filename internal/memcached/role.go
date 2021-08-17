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
func (t *RBAC) IndexRole(ctx context.Context, role internal.Roles) error {
	return t.orig.IndexRole(ctx, role)
}

func (t *RBAC) GetRole(ctx context.Context, roleId string) (internal.Roles, error) {
	key := "role_" + roleId
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			res, err := t.orig.GetRole(ctx, roleId)
			if err != nil {
				return internal.Roles{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetRole")
			}
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
		return internal.Roles{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.Roles
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.Roles{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}
func (t *RBAC) DeleteRole(ctx context.Context, roleId string) error {
	//get all accountrole
	ar, err := t.orig.AccountRoleByRoleReturnId(ctx, roleId)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetAccountRoleByRoleReturnId")
	}
	//delete first the role
	err = t.orig.DeleteRole(ctx, roleId)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.DeleteRole")
	}
	//get all taskrole
	rt, err := t.orig.RoleTaskByRoleReturnId(ctx, roleId)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetRoleTaskByRoleReturnId")
	}
	for _, value := range ar {
		err = t.orig.DeleteAccountRole(ctx, value)
		if err != nil {
			return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.DeleteAccountRole")
		}
	}
	for _, value := range rt {
		err = t.orig.DeleteRoleTask(ctx, value)
		if err != nil {
			return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.DeleteRoleTask")
		}
	}
	return nil
}

func (t *RBAC) ListRole(ctx context.Context, args internal.ListArgs) (internal.ListRole, error) {
	key := newKey("listrole", args)
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			listacc, err := t.orig.ListRole(ctx, args)
			if err != nil {
				return internal.ListRole{}, err
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
		return internal.ListRole{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.ListRole
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.ListRole{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}

func (t *RBAC) UpdateRole(ctx context.Context, role internal.Roles) error {
	err := t.orig.DeleteRole(ctx, role.Id)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.DeleteRole")
	}
	return t.orig.IndexRole(ctx, role)
}
