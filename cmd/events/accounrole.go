package events

import (
	"context"
	"rbac/internal"
)

func (r *RBACEvents) AccountRoleCreated(accountRole internal.AccountRoles) error {
	if err := r.cache.IndexAccountRole(context.Background(), accountRole); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) AccountRoleDeleted(id string) error {
	if err := r.cache.DeleteAccountRole(context.Background(), id); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) AccountRoleUpdated(accountRole internal.AccountRoles) error {
	if err := r.cache.UpdateAccountRole(context.Background(), accountRole); err != nil {
		return err
	}
	return nil
}
