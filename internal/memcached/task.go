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
func (t *RBAC) IndexTask(ctx context.Context, task internal.Tasks) error {
	return t.orig.IndexTask(ctx, task)
}

func (t *RBAC) GetTask(ctx context.Context, taskId string) (internal.Tasks, error) {
	key := "task_" + taskId
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			res, err := t.orig.GetTask(ctx, taskId)
			if err != nil {
				return internal.Tasks{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetTask")
			}
			res.HelpText, err = t.GetHelpTextByTask(ctx, taskId)
			if err != nil {
				return internal.Tasks{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetHelpText")
			}
			res.Menu, err = t.GetMenuByTask(ctx, taskId)
			if err != nil {
				return internal.Tasks{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetMenu")
			}
			res.Navigation, err = t.GetNavigationByTask(ctx, taskId)
			if err != nil {
				return internal.Tasks{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetNavigation")
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
		return internal.Tasks{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.Tasks
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.Tasks{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}
func (t *RBAC) DeleteTask(ctx context.Context, taskId string) error {
	return t.orig.DeleteTask(ctx, taskId)
}

func (t *RBAC) ListTask(ctx context.Context, args internal.ListArgs) (internal.ListTask, error) {
	key := newKey("listtask", args)
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			listacc, err := t.orig.ListTask(ctx, args)
			if err != nil {
				return internal.ListTask{}, err
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
		return internal.ListTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.ListTask
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.ListTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}
