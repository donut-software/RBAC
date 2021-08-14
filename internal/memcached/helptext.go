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
func (t *RBAC) IndexHelpText(ctx context.Context, helptext internal.HelpText) error {
	return t.orig.IndexHelpText(ctx, helptext)
}

func (t *RBAC) GetHelpText(ctx context.Context, helptextId string) (internal.HelpText, error) {
	key := "helptext_" + helptextId
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			res, err := t.orig.GetHelpText(ctx, helptextId)
			if err != nil {
				return internal.HelpText{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetHelpText")
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
		return internal.HelpText{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.HelpText
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.HelpText{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}
func (t *RBAC) DeleteHelpText(ctx context.Context, roleId string) error {
	return t.orig.DeleteHelpText(ctx, roleId)
}

func (t *RBAC) GetHelpTextByTask(ctx context.Context, taskid string) (internal.HelpText, error) {
	key := "helptextbytask_" + taskid
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			res, err := t.orig.HelpTextByTask(ctx, taskid)
			if err != nil {
				return internal.HelpText{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.AccountRoleByRole")
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
		return internal.HelpText{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.HelpText
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.HelpText{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}

func (t *RBAC) ListHelpText(ctx context.Context, args internal.ListArgs) (internal.ListHelpText, error) {
	key := newKey("listhelptext", args)
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			listacc, err := t.orig.ListHelpText(ctx, args)
			if err != nil {
				return internal.ListHelpText{}, err
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
		return internal.ListHelpText{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.ListHelpText
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.ListHelpText{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}
