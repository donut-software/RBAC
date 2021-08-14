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
func (t *RBAC) IndexNavigation(ctx context.Context, menu internal.Navigation) error {
	return t.orig.IndexNavigation(ctx, menu)
}

func (t *RBAC) GetNavigation(ctx context.Context, navigationId string) (internal.Navigation, error) {
	key := "navigation_" + navigationId
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			res, err := t.orig.GetNavigation(ctx, navigationId)
			if err != nil {
				return internal.Navigation{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetNavigation")
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
		return internal.Navigation{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.Navigation
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.Navigation{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}
func (t *RBAC) DeleteNavigation(ctx context.Context, roleId string) error {
	return t.orig.DeleteNavigation(ctx, roleId)
}

func (t *RBAC) GetNavigationByTask(ctx context.Context, taskid string) (internal.NavigationByTask, error) {
	key := "menubytask_" + taskid
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			res, err := t.orig.NavigationByTask(ctx, &taskid)
			if err != nil {
				return internal.NavigationByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.AccountRoleByRole")
			}
			task, err := t.orig.GetTask(ctx, res.Task.Id)
			if err != nil {
				return internal.NavigationByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetRole")
			}
			res.Task = task

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
		return internal.NavigationByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.NavigationByTask
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.NavigationByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}

func (t *RBAC) ListNavigation(ctx context.Context, args internal.ListArgs) (internal.ListNavigation, error) {
	key := newKey("listnavigation", args)
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			listacc, err := t.orig.ListNavigation(ctx, args)
			if err != nil {
				return internal.ListNavigation{}, err
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
		return internal.ListNavigation{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.ListNavigation
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.ListNavigation{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}
