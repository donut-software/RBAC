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
func (t *RBAC) IndexMenu(ctx context.Context, menu internal.Menu) error {
	return t.orig.IndexMenu(ctx, menu)
}

func (t *RBAC) GetMenu(ctx context.Context, menuId string) (internal.Menu, error) {
	key := "menu_" + menuId
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			res, err := t.orig.GetMenu(ctx, menuId)
			if err != nil {
				return internal.Menu{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetMenu")
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
		return internal.Menu{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.Menu
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.Menu{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}
func (t *RBAC) DeleteMenu(ctx context.Context, roleId string) error {
	return t.orig.DeleteMenu(ctx, roleId)
}

func (t *RBAC) GetMenuByTask(ctx context.Context, taskid string) ([]internal.Menu, error) {
	key := "menubytask_" + taskid
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			res, err := t.orig.MenuByTask(ctx, taskid)
			if err != nil {
				return []internal.Menu{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.AccountRoleByRole")
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
		return []internal.Menu{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res []internal.Menu
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return []internal.Menu{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}

func (t *RBAC) ListMenu(ctx context.Context, args internal.ListArgs) (internal.ListMenu, error) {
	key := newKey("listmenu", args)
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			listacc, err := t.orig.ListMenu(ctx, args)
			if err != nil {
				return internal.ListMenu{}, err
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
		return internal.ListMenu{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.ListMenu
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.ListMenu{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}
