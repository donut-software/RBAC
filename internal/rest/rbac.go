package rest

import (
	"context"
	"net/http"
	"rbac/internal"
	"rbac/internal/tokenmaker"

	"github.com/gorilla/mux"
)

//go:generate counterfeiter -o resttesting/rbac_service.gen.go . RBACService
type RBACService interface {
	Logout(ctx context.Context) error
	Login(ctx context.Context, username string, password string) error
	CreateAccount(ctx context.Context, account internal.Account, password string) (string, error)
	Account(ctx context.Context, username string) (internal.Account, error)
	AccountByID(ctx context.Context, id string) (internal.Account, error)
	UpdateProfile(ctx context.Context, profile internal.Profile) error
	ChangePassword(ctx context.Context, username string, password string) error
	DeleteAccount(ctx context.Context, username string) error
	ListAccount(ctx context.Context, args internal.ListArgs) (internal.ListAccount, error)
	IsAllowed(ctx context.Context, username string, task string) (bool, error)

	CreateRole(ctx context.Context, rolename string) (string, error)
	Role(ctx context.Context, id string) (internal.Roles, error)
	UpdateRole(ctx context.Context, id string, rolename string) error
	ListRole(ctx context.Context, args internal.ListArgs) (internal.ListRole, error)
	DeleteRole(ctx context.Context, id string) error

	CreateAccountRole(ctx context.Context, accountid string, roleid string) error
	AccountRole(ctx context.Context, accountRoleId string) (internal.AccountRoles, error)
	UpdateAccountRole(ctx context.Context, accountId string, roleId string, id string) error
	ListAccountRole(ctx context.Context, args internal.ListArgs) (internal.ListAccountRole, error)
	AccountRoleByAccount(ctx context.Context, accountRoleId string) (internal.AccountRoleByAccountResult, error)
	AccountRoleByRole(ctx context.Context, id string) (internal.AccountRoleByRoleResult, error)
	DeleteAccountRole(ctx context.Context, id string) error

	CreateTask(ctx context.Context, taskname string) (string, error)
	Task(ctx context.Context, id string) (internal.Tasks, error)
	UpdateTask(ctx context.Context, id string, taskname string) error
	ListTask(ctx context.Context, args internal.ListArgs) (internal.ListTask, error)
	DeleteTask(ctx context.Context, id string) error

	CreateRoleTask(ctx context.Context, taskid string, roleid string) error
	RoleTask(ctx context.Context, roleTaskId string) (internal.RoleTasks, error)
	UpdateRoleTask(ctx context.Context, taskId string, roleId string, id string) error
	ListRoleTask(ctx context.Context, args internal.ListArgs) (internal.ListRoleTask, error)
	DeleteRoleTask(ctx context.Context, id string) error

	CreateHelpText(ctx context.Context, helptext internal.HelpText) error
	HelpText(ctx context.Context, id string) (internal.HelpText, error)
	UpdateHelpText(ctx context.Context, helptext internal.HelpText) error
	DeleteHelpText(ctx context.Context, id string) error
	ListHelpText(ctx context.Context, args internal.ListArgs) (internal.ListHelpText, error)

	CreateMenu(ctx context.Context, menu internal.Menu) error
	Menu(ctx context.Context, id string) (internal.Menu, error)
	UpdateMenu(ctx context.Context, menu internal.Menu) error
	DeleteMenu(ctx context.Context, id string) error
	ListMenu(ctx context.Context, args internal.ListArgs) (internal.ListMenu, error)

	CreateNavigation(ctx context.Context, navigation internal.Navigation) error
	Navigation(ctx context.Context, id string) (internal.Navigation, error)
	UpdateNavigation(ctx context.Context, navigation internal.Navigation) error
	ListNavigation(ctx context.Context, args internal.ListArgs) (internal.ListNavigation, error)
	DeleteNavigation(ctx context.Context, id string) error

	CreateToken(username string) (string, error)
	VerifyToken(token string) (*tokenmaker.Payload, error)
}

type RBACHandler struct {
	svc RBACService
}

func NewRBACHandler(svc RBACService) *RBACHandler {
	return &RBACHandler{
		svc: svc,
	}
}

func (rb *RBACHandler) Register(r *mux.Router) {

	v0 := r.PathPrefix("/v0/").Subrouter()

	r.HandleFunc("/v0/login", rb.login).Methods(http.MethodPost)
	r.HandleFunc("/v0/register", rb.register).Methods(http.MethodPost)

	v0.Use(rb.middleware)
	accountRouter := v0.PathPrefix("/accounts/").Subrouter()
	accountRouter.HandleFunc("/logout", rb.logout).Methods(http.MethodPost)
	accountRouter.HandleFunc("/me", rb.me).Methods(http.MethodGet)

	accountRouter.HandleFunc("/{username}", rb.account).Methods(http.MethodGet)
	accountRouter.HandleFunc("/roles/{username}", rb.getAccountRoleByAccount).Methods(http.MethodGet)
	accountRouter.HandleFunc("/", rb.listaccount).Methods(http.MethodGet)
	accountRouter.HandleFunc("/", rb.updateProfile).Methods(http.MethodPut)
	accountRouter.HandleFunc("/{username}", rb.deleteAccount).Methods(http.MethodDelete)

	roleRouter := v0.PathPrefix("/roles/").Subrouter()
	roleRouter.HandleFunc("/", rb.createRole).Methods(http.MethodPost)
	roleRouter.HandleFunc("/{roleId}", rb.role).Methods(http.MethodGet)
	roleRouter.HandleFunc("/accounts/{roleId}", rb.getAccountRoleByRole).Methods(http.MethodGet)
	roleRouter.HandleFunc("/", rb.updateRole).Methods(http.MethodPut)
	roleRouter.HandleFunc("/", rb.listrole).Methods(http.MethodGet)
	roleRouter.HandleFunc("/{roleId}", rb.deleteRole).Methods(http.MethodDelete)

	accountroleRouter := v0.PathPrefix("/accountroles/").Subrouter()
	accountroleRouter.HandleFunc("/", rb.createAccountRole).Methods(http.MethodPost)
	accountroleRouter.HandleFunc("/{accountRoleId}", rb.accountRole).Methods(http.MethodGet)
	accountroleRouter.HandleFunc("/", rb.updateAccountRole).Methods(http.MethodPut)
	accountroleRouter.HandleFunc("/", rb.listAccountRole).Methods(http.MethodGet)
	accountroleRouter.HandleFunc("/{accountRoleId}", rb.deleteAccountRole).Methods(http.MethodDelete)

	taskRouter := v0.PathPrefix("/task/").Subrouter()
	taskRouter.HandleFunc("/", rb.createTask).Methods(http.MethodPost)
	taskRouter.HandleFunc("/{taskId}", rb.task).Methods(http.MethodGet)
	taskRouter.HandleFunc("/", rb.updateTask).Methods(http.MethodPut)
	taskRouter.HandleFunc("/", rb.listtask).Methods(http.MethodGet)
	taskRouter.HandleFunc("/{taskId}", rb.deleteTask).Methods(http.MethodDelete)

	roletaskRouter := v0.PathPrefix("/roletask/").Subrouter()
	roletaskRouter.HandleFunc("/", rb.createRoleTask).Methods(http.MethodPost)
	roletaskRouter.HandleFunc("/{roleTaskId}", rb.roleTask).Methods(http.MethodGet)
	roletaskRouter.HandleFunc("/", rb.updateRoleTask).Methods(http.MethodPut)
	roletaskRouter.HandleFunc("/", rb.listRoleTask).Methods(http.MethodGet)
	roletaskRouter.HandleFunc("/{roleTaskId}", rb.deleteRoleTask).Methods(http.MethodDelete)

	helptextRouter := v0.PathPrefix("/helptext/").Subrouter()
	helptextRouter.HandleFunc("/", rb.createHelpText).Methods(http.MethodPost)
	helptextRouter.HandleFunc("/{helpTextId}", rb.helpText).Methods(http.MethodGet)
	helptextRouter.HandleFunc("/", rb.updateHelpText).Methods(http.MethodPut)
	helptextRouter.HandleFunc("/", rb.listHelpText).Methods(http.MethodGet)
	helptextRouter.HandleFunc("/{helpTextId}", rb.deleteHelpText).Methods(http.MethodDelete)

	menuRouter := v0.PathPrefix("/menu/").Subrouter()
	menuRouter.HandleFunc("/", rb.createMenu).Methods(http.MethodPost)
	menuRouter.HandleFunc("/{menuId}", rb.menu).Methods(http.MethodGet)
	menuRouter.HandleFunc("/", rb.updateMenu).Methods(http.MethodPut)
	menuRouter.HandleFunc("/", rb.listMenu).Methods(http.MethodGet)
	menuRouter.HandleFunc("/{menuId}", rb.deleteMenu).Methods(http.MethodDelete)

	navigationRouter := v0.PathPrefix("/navigation/").Subrouter()
	navigationRouter.HandleFunc("/", rb.createNavigation).Methods(http.MethodPost)
	navigationRouter.HandleFunc("/{navigationId}", rb.navigation).Methods(http.MethodGet)
	navigationRouter.HandleFunc("/", rb.updateNavigation).Methods(http.MethodPut)
	navigationRouter.HandleFunc("/", rb.listNavigation).Methods(http.MethodGet)
	navigationRouter.HandleFunc("/{navigationId}", rb.deleteNavigation).Methods(http.MethodDelete)

}
