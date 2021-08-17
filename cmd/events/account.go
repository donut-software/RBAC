package events

import (
	"context"
	"rbac/internal"
)

func (r *RBACEvents) AccountCreated(account internal.Account) error {
	if err := r.cache.IndexAccount(context.Background(), account); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) AccountDeleted(id string) error {
	if err := r.cache.DeleteAccount(context.Background(), id); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) AccountUpdated(profile internal.Profile) error {
	if err := r.cache.UpdateProfile(context.Background(), profile); err != nil {
		return err
	}
	return nil
}
