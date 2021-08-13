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
func (t *RBAC) IndexRoleTask(ctx context.Context, roletask internal.RoleTasks) error {
	return t.orig.IndexRoleTask(ctx, roletask)
}

func (t *RBAC) GetRoleTask(ctx context.Context, roletaskId string) (internal.RoleTasks, error) {
	key := "roletask_" + roletaskId
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			res, err := t.orig.GetRoleTask(ctx, roletaskId)
			if err != nil {
				return internal.RoleTasks{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetRoleTask")
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
		return internal.RoleTasks{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.RoleTasks
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.RoleTasks{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}
func (t *RBAC) DeleteRoleTask(ctx context.Context, roletaskId string) error {
	return t.orig.DeleteRoleTask(ctx, roletaskId)
}

func (t *RBAC) GetRoleTaskByRole(ctx context.Context, roleid string) (internal.RoleTaskByRole, error) {
	key := "roletaskbyrole_" + roleid
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			res, err := t.orig.RoleTaskByRole(ctx, &roleid)
			if err != nil {
				return internal.RoleTaskByRole{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.RoleTaskByRole")
			}
			role, err := t.orig.GetRole(ctx, res.Role.Id)
			if err != nil {
				return internal.RoleTaskByRole{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetRole")
			}
			res.Role = role
			tasks := []internal.Tasks{}
			for _, value := range res.Tasks {
				task, err := t.orig.GetTask(ctx, value.Id)
				if err != nil {
					return internal.RoleTaskByRole{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetTask")
				}
				tasks = append(tasks, task)
			}
			res.Tasks = tasks
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
		return internal.RoleTaskByRole{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.RoleTaskByRole
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.RoleTaskByRole{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}

func (t *RBAC) GetRoleTaskByTask(ctx context.Context, taskid string) (internal.RoleTaskByTask, error) {
	key := "roletaskbytask_" + taskid
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			res, err := t.orig.RoleTaskByTask(ctx, &taskid)
			if err != nil {
				return internal.RoleTaskByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.RoleTaskByTask")
			}
			task, err := t.orig.GetTask(ctx, res.Task.Id)
			if err != nil {
				return internal.RoleTaskByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetTask")
			}
			res.Task = task
			roles := []internal.Roles{}
			for _, value := range res.Roles {
				role, err := t.orig.GetRole(ctx, value.Id)
				if err != nil {
					return internal.RoleTaskByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetRole")
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
		return internal.RoleTaskByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.RoleTaskByTask
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.RoleTaskByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}
