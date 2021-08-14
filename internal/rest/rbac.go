package rest

import (
	"context"
	"net/http"
	"rbac/internal"

	"github.com/gorilla/mux"
)

//go:generate counterfeiter -o resttesting/rbac_service.gen.go . RBACService
type RBACService interface {
	CreateAccount(ctx context.Context, account internal.Account, password string) error
	Account(ctx context.Context, username string) (internal.Account, error)
	AccountByID(ctx context.Context, id string) (internal.Account, error)
	UpdateProfile(ctx context.Context, profile internal.Profile) error
	ChangePassword(ctx context.Context, username string, password string) error
	DeleteAccount(ctx context.Context, username string) error
	ListAccount(ctx context.Context, args internal.ListArgs) (internal.ListAccount, error)

	CreateRole(ctx context.Context, rolename string) error
	Role(ctx context.Context, id string) (internal.Roles, error)
	UpdateRole(ctx context.Context, id string, rolename string) error
	ListRole(ctx context.Context, args internal.ListArgs) (internal.ListRole, error)

	CreateAccountRole(ctx context.Context, accountid string, roleid string) error
	AccountRole(ctx context.Context, accountRoleId string) (internal.AccountRoles, error)
	UpdateAccountRole(ctx context.Context, accountId string, roleId string, id string) error
	ListAccountRole(ctx context.Context, args internal.ListArgs) (internal.ListAccountRole, error)

	CreateTask(ctx context.Context, taskname string) error
	Task(ctx context.Context, id string) (internal.Tasks, error)
	UpdateTask(ctx context.Context, id string, taskname string) error

	CreateRoleTask(ctx context.Context, taskid string, roleid string) error
	RoleTask(ctx context.Context, roleTaskId string) (internal.RoleTasks, error)
	UpdateRoleTask(ctx context.Context, taskId string, roleId string, id string) error

	CreateHelpText(ctx context.Context, helptext internal.HelpText) error
	HelpText(ctx context.Context, id string) (internal.HelpText, error)
	UpdateHelpText(ctx context.Context, helptext internal.HelpText) error

	CreateMenu(ctx context.Context, menu internal.Menu) error
	Menu(ctx context.Context, id string) (internal.Menu, error)
	UpdateMenu(ctx context.Context, menu internal.Menu) error

	CreateNavigation(ctx context.Context, navigation internal.Navigation) error
	Navigation(ctx context.Context, id string) (internal.Navigation, error)
	UpdateNavigation(ctx context.Context, navigation internal.Navigation) error
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

	accountRouter := r.PathPrefix("/accounts/").Subrouter()
	accountRouter.HandleFunc("/register", rb.register).Methods(http.MethodPost)
	accountRouter.HandleFunc("/{username}", rb.account).Methods(http.MethodGet)
	accountRouter.HandleFunc("/", rb.listaccount).Methods(http.MethodGet)
	accountRouter.HandleFunc("/", rb.updateProfile).Methods(http.MethodPut)
	accountRouter.HandleFunc("/{username}", rb.deleteAccount).Methods(http.MethodDelete)

	roleRouter := r.PathPrefix("/roles/").Subrouter()
	roleRouter.HandleFunc("/", rb.createRole).Methods(http.MethodPost)
	roleRouter.HandleFunc("/{roleId}", rb.role).Methods(http.MethodGet)
	roleRouter.HandleFunc("/", rb.updateRole).Methods(http.MethodPut)
	roleRouter.HandleFunc("/", rb.listrole).Methods(http.MethodGet)

	accountroleRouter := r.PathPrefix("/accountroles/").Subrouter()
	accountroleRouter.HandleFunc("/", rb.createAccountRole).Methods(http.MethodPost)
	accountroleRouter.HandleFunc("/{accountRoleId}", rb.accountRole).Methods(http.MethodGet)
	accountroleRouter.HandleFunc("/", rb.updateAccountRole).Methods(http.MethodPut)
	accountroleRouter.HandleFunc("/", rb.listAccountRole).Methods(http.MethodGet)

	taskRouter := r.PathPrefix("/task/").Subrouter()
	taskRouter.HandleFunc("/", rb.createTask).Methods(http.MethodPost)
	taskRouter.HandleFunc("/{taskId}", rb.task).Methods(http.MethodGet)
	taskRouter.HandleFunc("/", rb.updateTask).Methods(http.MethodPut)

	roletaskRouter := r.PathPrefix("/roletask/").Subrouter()
	roletaskRouter.HandleFunc("/", rb.createRoleTask).Methods(http.MethodPost)
	roletaskRouter.HandleFunc("/{roleTaskId}", rb.roleTask).Methods(http.MethodGet)
	roletaskRouter.HandleFunc("/", rb.updateRoleTask).Methods(http.MethodPut)

	helptextRouter := r.PathPrefix("/helptext/").Subrouter()
	helptextRouter.HandleFunc("/", rb.createHelpText).Methods(http.MethodPost)
	helptextRouter.HandleFunc("/{helpTextId}", rb.helpText).Methods(http.MethodGet)
	helptextRouter.HandleFunc("/", rb.updateHelpText).Methods(http.MethodPut)

	menuRouter := r.PathPrefix("/menu/").Subrouter()
	menuRouter.HandleFunc("/", rb.createMenu).Methods(http.MethodPost)
	menuRouter.HandleFunc("/{menuId}", rb.menu).Methods(http.MethodGet)
	menuRouter.HandleFunc("/", rb.updateMenu).Methods(http.MethodPut)

	navigationRouter := r.PathPrefix("/navigation/").Subrouter()
	navigationRouter.HandleFunc("/", rb.createNavigation).Methods(http.MethodPost)
	navigationRouter.HandleFunc("/{menuId}", rb.navigation).Methods(http.MethodGet)
	navigationRouter.HandleFunc("/", rb.updateNavigation).Methods(http.MethodPut)

}
