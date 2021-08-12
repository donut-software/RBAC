package service

import (
	"context"
	"fmt"
	"rbac/internal"
)

func (r *RBAC) CreateRoleTask(ctx context.Context, taskid string, roleid string) error {
	err := r.repo.CreateRoleTasks(ctx, taskid, roleid)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (r *RBAC) RoleTask(ctx context.Context, roleTaskId string) (internal.RoleTasks, error) {
	role, err := r.repo.RoleTask(ctx, roleTaskId)
	if err != nil {
		fmt.Println(err)
		return internal.RoleTasks{}, err
	}
	return role, err
}
func (r *RBAC) UpdateRoleTask(ctx context.Context, taskId string, roleId string, id string) error {
	err := r.repo.UpdateRoleTask(ctx, taskId, roleId, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
func (r *RBAC) DeleteRoleTask(ctx context.Context, id string) error {
	err := r.repo.DeleteRoleTask(ctx, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
