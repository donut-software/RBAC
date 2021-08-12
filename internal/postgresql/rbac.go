package postgresql

import (
	"context"
	"database/sql"
	"rbac/internal"
)

type RBAC interface {
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
}

func NewRBAC(db *sql.DB) RBAC {
	return &Store{
		db: db,
		q:  New(db),
	}
}
