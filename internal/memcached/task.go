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
func (t *RBAC) IndexTask(ctx context.Context, role internal.Tasks) error {
	return t.orig.IndexTask(ctx, role)
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
func (t *RBAC) DeleteTask(ctx context.Context, roleId string) error {
	return t.orig.DeleteTask(ctx, roleId)
}
