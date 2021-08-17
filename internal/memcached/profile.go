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
func (t *RBAC) IndexProfile(ctx context.Context, profile internal.Profile) error {
	return t.orig.IndexProfile(ctx, profile)
}

func (t *RBAC) GetProfile(ctx context.Context, profileid string) (internal.Profile, error) {
	key := "profile_" + profileid
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			res, err := t.orig.GetProfile(ctx, profileid)
			if err != nil {
				return internal.Profile{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetProfile")
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
		return internal.Profile{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.Profile
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.Profile{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}
func (t *RBAC) DeleteProfile(ctx context.Context, roleId string) error {
	return t.orig.DeleteProfile(ctx, roleId)
}

func (t *RBAC) UpdateProfile(ctx context.Context, profile internal.Profile) error {
	err := t.orig.DeleteProfile(ctx, profile.Id)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.DelerteProfie")
	}
	return t.orig.IndexProfile(ctx, profile)
}
