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
	return t.orig.DeleteRole(ctx, roleId)
}
