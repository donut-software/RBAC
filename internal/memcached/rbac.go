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

type Datastore interface {
	IndexAccount(ctx context.Context, account internal.Account) error
	GetAccount(ctx context.Context, username string) (internal.Account, error)
	DeleteAccount(ctx context.Context, username string) error

	IndexProfile(ctx context.Context, profile internal.Profile) error
	GetProfile(ctx context.Context, profileId string) (internal.Profile, error)
	DeleteProfile(ctx context.Context, profileId string) error

	IndexRole(ctx context.Context, role internal.Roles) error
	DeleteRole(ctx context.Context, roleId string) error
	GetRole(ctx context.Context, roleId string) (internal.Roles, error)

	IndexAccountRole(ctx context.Context, accRole internal.AccountRoles) error
	GetAccountRole(ctx context.Context, accRoleId string) (internal.AccountRoles, error)
	DeleteAccountRole(ctx context.Context, accRoleId string) error

	IndexTask(ctx context.Context, task internal.Tasks) error
	DeleteTask(ctx context.Context, taskId string) error
	GetTask(ctx context.Context, taskId string) (internal.Tasks, error)

	IndexRoleTask(ctx context.Context, roletask internal.RoleTasks) error
	DeleteRoleTask(ctx context.Context, roletaskId string) error
	GetRoleTask(ctx context.Context, roletaskId string) (internal.RoleTasks, error)

	IndexHelpText(ctx context.Context, helptext internal.HelpText) error
	DeleteHelpText(ctx context.Context, helptextId string) error
	GetHelpText(ctx context.Context, helptextId string) (internal.HelpText, error)

	IndexMenu(ctx context.Context, menu internal.Menu) error
	DeleteMenu(ctx context.Context, menuId string) error
	GetMenu(ctx context.Context, menuId string) (internal.Menu, error)

	IndexNavigation(ctx context.Context, navigation internal.Navigation) error
	DeleteNavigation(ctx context.Context, navigationId string) error
	GetNavigation(ctx context.Context, navigationId string) (internal.Navigation, error)
}
type RBAC struct {
	client *memcache.Client
	orig   Datastore
	logger *zap.Logger
}

func NewRBAC(client *memcache.Client, orig Datastore, logger *zap.Logger) *RBAC {
	return &RBAC{
		client: client,
		orig:   orig,
		logger: logger,
	}
}

// Index ...
func (t *RBAC) IndexAccount(ctx context.Context, account internal.Account) error {
	err := t.orig.IndexProfile(ctx, account.Profile)
	if err != nil {
		return err
	}
	return t.orig.IndexAccount(ctx, account)
}

func (t *RBAC) GetAccount(ctx context.Context, username string) (internal.Account, error) {
	key := username
	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))
			res, err := t.orig.GetAccount(ctx, username)
			if err != nil {
				return internal.Account{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetAccount")
			}
			prof, err := t.orig.GetProfile(ctx, res.Profile.Id)
			if err != nil {
				return internal.Account{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.GetProfile")
			}
			res.Profile = prof
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
		return internal.Account{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}
	t.logger.Info("values found", zap.String("key", string(key)))
	var res internal.Account
	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.Account{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}
	return res, nil
}
