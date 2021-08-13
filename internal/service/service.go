package service

import (
	"rbac/internal"

	"golang.org/x/net/context"
)

type RBACRepository interface {
	CreateAccount(ctx context.Context, account internal.Account, password string) error
	Account(ctx context.Context, username string) (internal.Account, error)
	UpdateProfile(ctx context.Context, profile internal.Profile) error
	ChangePassword(ctx context.Context, username string, password string) error
	DeleteAccount(ctx context.Context, username string) error

	CreateRole(ctx context.Context, rolename string) error
	Role(ctx context.Context, id string) (internal.Roles, error)
	UpdateRole(ctx context.Context, id string, rolename string) error
	DeleteRole(ctx context.Context, id string) error

	CreateAccountRole(ctx context.Context, accountid string, roleid string) error
	AccountRole(ctx context.Context, accountRoleId string) (internal.AccountRoles, error)
	UpdateAccountRole(ctx context.Context, accountId string, roleId string, id string) error
	DeleteAccountRole(ctx context.Context, id string) error

	CreateTask(ctx context.Context, taskname string) error
	Task(ctx context.Context, id string) (internal.Tasks, error)
	UpdateTask(ctx context.Context, id string, taskname string) error
	DeleteTask(ctx context.Context, id string) error

	CreateRoleTasks(ctx context.Context, taskid string, roleid string) error
	RoleTask(ctx context.Context, roleTaskId string) (internal.RoleTasks, error)
	UpdateRoleTask(ctx context.Context, taskId string, roleId string, id string) error
	DeleteRoleTask(ctx context.Context, id string) error

	CreateHelpText(ctx context.Context, helptext internal.HelpText) error
	HelpText(ctx context.Context, id string) (internal.HelpText, error)
	UpdateHelpText(ctx context.Context, helptext internal.HelpText) error
	DeleteHelpText(ctx context.Context, id string) error

	CreateMenu(ctx context.Context, menu internal.Menu) error
	Menu(ctx context.Context, id string) (internal.Menu, error)
	UpdateMenu(ctx context.Context, menu internal.Menu) error
	DeleteMenu(ctx context.Context, id string) error

	CreateNavigation(ctx context.Context, menu internal.Navigation) error
	Navigation(ctx context.Context, id string) (internal.Navigation, error)
	UpdateNavigation(ctx context.Context, menu internal.Navigation) error
	DeleteNavigation(ctx context.Context, id string) error
}

type RBACSearchRepository interface {
	IndexAccount(ctx context.Context, account internal.Account) error
	// DeleteAccount(ctx context.Context, username string) error
	GetAccount(ctx context.Context, username string) (internal.Account, error)
	ListAccount(ctx context.Context, args internal.ListAccountArgs) (internal.ListAccount, error)
	// IndexProfile(ctx context.Context, profile internal.Profile) error
}
type RBAC struct {
	repo   RBACRepository
	search RBACSearchRepository
}

func NewRBAC(repo RBACRepository, search RBACSearchRepository) *RBAC {
	return &RBAC{
		repo:   repo,
		search: search,
	}
}
