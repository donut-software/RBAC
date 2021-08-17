package events

import (
	"context"
	"rbac/internal"
)

func (r *RBACEvents) RoleCreated(role internal.Roles) error {
	if err := r.cache.IndexRole(context.Background(), role); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) RoleDeleted(id string) error {
	if err := r.cache.DeleteRole(context.Background(), id); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) RoleUpdated(role internal.Roles) error {
	if err := r.cache.UpdateRole(context.Background(), role); err != nil {
		return err
	}
	return nil
}
