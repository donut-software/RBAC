package memcached

import (
	"context"
	"rbac/internal"

	"github.com/bradfitz/gomemcache/memcache"
	"go.uber.org/zap"
)

type Datastore interface {
	IndexAccount(ctx context.Context, account internal.Account) error
	GetAccount(ctx context.Context, username string) (internal.Account, error)
	GetAccountById(ctx context.Context, id *string) (internal.Account, error)
	DeleteAccount(ctx context.Context, username string) error
	ListAccount(ctx context.Context, args internal.ListArgs) (internal.ListAccount, error)

	IndexProfile(ctx context.Context, profile internal.Profile) error
	GetProfile(ctx context.Context, profileId string) (internal.Profile, error)
	DeleteProfile(ctx context.Context, profileId string) error

	IndexRole(ctx context.Context, role internal.Roles) error
	DeleteRole(ctx context.Context, roleId string) error
	GetRole(ctx context.Context, roleId string) (internal.Roles, error)
	ListRole(ctx context.Context, args internal.ListArgs) (internal.ListRole, error)

	IndexAccountRole(ctx context.Context, accRole internal.AccountRoles) error
	GetAccountRole(ctx context.Context, accRoleId string) (internal.AccountRoles, error)
	DeleteAccountRole(ctx context.Context, accRoleId string) error
	AccountRoleByAccount(ctx context.Context, username *string) (internal.AccountRoleByAccountResult, error)
	AccountRoleByRole(ctx context.Context, roleId *string) (internal.AccountRoleByRoleResult, error)
	ListAccountRole(ctx context.Context, args internal.ListArgs) (internal.ListAccountRole, error)

	IndexTask(ctx context.Context, task internal.Tasks) error
	DeleteTask(ctx context.Context, taskId string) error
	GetTask(ctx context.Context, taskId string) (internal.Tasks, error)
	ListTask(ctx context.Context, args internal.ListArgs) (internal.ListTask, error)

	IndexRoleTask(ctx context.Context, roletask internal.RoleTasks) error
	DeleteRoleTask(ctx context.Context, roletaskId string) error
	GetRoleTask(ctx context.Context, roletaskId string) (internal.RoleTasks, error)
	RoleTaskByRole(ctx context.Context, roleId *string) (internal.RoleTaskByRole, error)
	RoleTaskByTask(ctx context.Context, taskId *string) (internal.RoleTaskByTask, error)
	ListRoleTask(ctx context.Context, args internal.ListArgs) (internal.ListRoleTask, error)

	IndexHelpText(ctx context.Context, helptext internal.HelpText) error
	DeleteHelpText(ctx context.Context, helptextId string) error
	GetHelpText(ctx context.Context, helptextId string) (internal.HelpText, error)
	HelpTextByTask(ctx context.Context, taskId *string) (internal.HelpTextByTask, error)
	ListHelpText(ctx context.Context, args internal.ListArgs) (internal.ListHelpText, error)

	IndexMenu(ctx context.Context, menu internal.Menu) error
	DeleteMenu(ctx context.Context, menuId string) error
	GetMenu(ctx context.Context, menuId string) (internal.Menu, error)
	MenuByTask(ctx context.Context, taskId *string) (internal.MenuByTask, error)
	ListMenu(ctx context.Context, args internal.ListArgs) (internal.ListMenu, error)

	IndexNavigation(ctx context.Context, navigation internal.Navigation) error
	DeleteNavigation(ctx context.Context, navigationId string) error
	GetNavigation(ctx context.Context, navigationId string) (internal.Navigation, error)
	NavigationByTask(ctx context.Context, taskId *string) (internal.NavigationByTask, error)
	ListNavigation(ctx context.Context, args internal.ListArgs) (internal.ListNavigation, error)
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
