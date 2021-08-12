package service

import (
	"context"
	"fmt"
	"rbac/internal"
)

func (r *RBAC) CreateRole(ctx context.Context, rolename string) error {
	err := r.repo.CreateRole(ctx, rolename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (r *RBAC) Role(ctx context.Context, id string) (internal.Roles, error) {
	role, err := r.repo.Role(ctx, id)
	if err != nil {
		fmt.Println(err)
		return internal.Roles{}, err
	}
	return role, err
}
func (r *RBAC) UpdateRole(ctx context.Context, id string, rolename string) error {
	err := r.repo.UpdateRole(ctx, id, rolename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
func (r *RBAC) DeleteRole(ctx context.Context, id string) error {
	err := r.repo.DeleteRole(ctx, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
