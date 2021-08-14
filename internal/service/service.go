package service

import (
	"rbac/internal"

	"golang.org/x/net/context"
)

type RBACRepository interface {
	Login(ctx context.Context, username string, password string) error
	CreateAccount(ctx context.Context, account internal.Account, password string) error
	Account(ctx context.Context, username string) (internal.Account, error)
	UpdateProfile(ctx context.Context, profile internal.Profile) error
	ChangePassword(ctx context.Context, username string, password string) error
	DeleteAccount(ctx context.Context, username string) error

	CreateRole(ctx context.Context, rolename string) (string, error)
	Role(ctx context.Context, id string) (internal.Roles, error)
	UpdateRole(ctx context.Context, id string, rolename string) error
	DeleteRole(ctx context.Context, id string) error

	CreateAccountRole(ctx context.Context, accountId string, roleId string) (string, error)
	AccountRole(ctx context.Context, accountRoleId string) (internal.AccountRoles, error)
	UpdateAccountRole(ctx context.Context, accountId string, roleId string, id string) error
	DeleteAccountRole(ctx context.Context, id string) error

	CreateTask(ctx context.Context, taskname string) (string, error)
	Task(ctx context.Context, id string) (internal.Tasks, error)
	UpdateTask(ctx context.Context, id string, taskname string) error
	DeleteTask(ctx context.Context, id string) error

	CreateRoleTasks(ctx context.Context, taskid string, roleid string) (string, error)
	RoleTask(ctx context.Context, roleTaskId string) (internal.RoleTasks, error)
	UpdateRoleTask(ctx context.Context, taskId string, roleId string, id string) error
	DeleteRoleTask(ctx context.Context, id string) error

	CreateHelpText(ctx context.Context, helptext internal.HelpText) (string, error)
	HelpText(ctx context.Context, id string) (internal.HelpText, error)
	UpdateHelpText(ctx context.Context, helptext internal.HelpText) error
	DeleteHelpText(ctx context.Context, id string) error

	CreateMenu(ctx context.Context, menu internal.Menu) (string, error)
	Menu(ctx context.Context, id string) (internal.Menu, error)
	UpdateMenu(ctx context.Context, menu internal.Menu) error
	DeleteMenu(ctx context.Context, id string) error

	CreateNavigation(ctx context.Context, menu internal.Navigation) (string, error)
	Navigation(ctx context.Context, id string) (internal.Navigation, error)
	UpdateNavigation(ctx context.Context, menu internal.Navigation) error
	DeleteNavigation(ctx context.Context, id string) error
}

type RBACSearchRepository interface {
	IndexAccount(ctx context.Context, account internal.Account) error
	GetAccount(ctx context.Context, username string) (internal.Account, error)
	GetAccountById(ctx context.Context, id string) (internal.Account, error)
	ListAccount(ctx context.Context, args internal.ListArgs) (internal.ListAccount, error)
	DeleteAccount(ctx context.Context, username string) error

	IndexProfile(ctx context.Context, profile internal.Profile) error
	GetProfile(ctx context.Context, profileid string) (internal.Profile, error)
	DeleteProfile(ctx context.Context, roleId string) error

	IndexRole(ctx context.Context, role internal.Roles) error
	GetRole(ctx context.Context, roleId string) (internal.Roles, error)
	DeleteRole(ctx context.Context, roleId string) error
	ListRole(ctx context.Context, args internal.ListArgs) (internal.ListRole, error)

	IndexAccountRole(ctx context.Context, accRole internal.AccountRoles) error
	GetAccountRole(ctx context.Context, accRoleId string) (internal.AccountRoles, error)
	GetAccountRoleByAccount(ctx context.Context, username string) (internal.AccountRoleByAccountResult, error)
	GetAccountRoleByRole(ctx context.Context, roleid string) (internal.AccountRoleByRoleResult, error)
	ListAccountRole(ctx context.Context, args internal.ListArgs) (internal.ListAccountRole, error)
	DeleteAccountRole(ctx context.Context, accRoleId string) error

	IndexTask(ctx context.Context, role internal.Tasks) error
	GetTask(ctx context.Context, taskId string) (internal.Tasks, error)
	DeleteTask(ctx context.Context, roleId string) error
	ListTask(ctx context.Context, args internal.ListArgs) (internal.ListTask, error)

	IndexRoleTask(ctx context.Context, roletask internal.RoleTasks) error
	GetRoleTask(ctx context.Context, roletaskId string) (internal.RoleTasks, error)
	DeleteRoleTask(ctx context.Context, roletaskId string) error
	GetRoleTaskByRole(ctx context.Context, roleid string) (internal.RoleTaskByRole, error)
	GetRoleTaskByTask(ctx context.Context, taskid string) (internal.RoleTaskByTask, error)
	ListRoleTask(ctx context.Context, args internal.ListArgs) (internal.ListRoleTask, error)

	IndexHelpText(ctx context.Context, helptext internal.HelpText) error
	GetHelpText(ctx context.Context, helptextId string) (internal.HelpText, error)
	GetHelpTextByTask(ctx context.Context, taskid string) (internal.HelpText, error)
	DeleteHelpText(ctx context.Context, roleId string) error
	ListHelpText(ctx context.Context, args internal.ListArgs) (internal.ListHelpText, error)

	IndexMenu(ctx context.Context, menu internal.Menu) error
	GetMenu(ctx context.Context, menuId string) (internal.Menu, error)
	DeleteMenu(ctx context.Context, roleId string) error
	GetMenuByTask(ctx context.Context, taskid string) ([]internal.Menu, error)
	ListMenu(ctx context.Context, args internal.ListArgs) (internal.ListMenu, error)

	IndexNavigation(ctx context.Context, menu internal.Navigation) error
	GetNavigation(ctx context.Context, navigationId string) (internal.Navigation, error)
	DeleteNavigation(ctx context.Context, roleId string) error
	GetNavigationByTask(ctx context.Context, taskid string) ([]internal.Navigation, error)
	ListNavigation(ctx context.Context, args internal.ListArgs) (internal.ListNavigation, error)
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
