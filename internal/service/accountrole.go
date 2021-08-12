package service

import (
	"context"
	"fmt"
	"rbac/internal"
)

func (r *RBAC) CreateAccountRole(ctx context.Context, accountid string, roleid string) error {
	err := r.repo.CreateAccountRole(ctx, accountid, roleid)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (r *RBAC) AccountRole(ctx context.Context, accountRoleId string) (internal.AccountRoles, error) {
	role, err := r.repo.AccountRole(ctx, accountRoleId)
	if err != nil {
		fmt.Println(err)
		return internal.AccountRoles{}, err
	}
	return role, err
}
func (r *RBAC) UpdateAccountRole(ctx context.Context, accountId string, roleId string, id string) error {
	err := r.repo.UpdateAccountRole(ctx, accountId, roleId, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
func (r *RBAC) DeleteAccountRole(ctx context.Context, id string) error {
	err := r.repo.DeleteAccountRole(ctx, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
