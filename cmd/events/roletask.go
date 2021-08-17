package events

import (
	"context"
	"rbac/internal"
)

func (r *RBACEvents) RoleTaskCreated(roleTask internal.RoleTasks) error {
	if err := r.cache.IndexRoleTask(context.Background(), roleTask); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) RoleTaskDeleted(id string) error {
	if err := r.cache.DeleteRoleTask(context.Background(), id); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) RoleTaskUpdated(roleTask internal.RoleTasks) error {
	if err := r.cache.UpdateRoleTask(context.Background(), roleTask); err != nil {
		return err
	}
	return nil
}
