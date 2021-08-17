package events

import (
	"context"
	"rbac/internal"
)

func (r *RBACEvents) TaskCreated(task internal.Tasks) error {
	if err := r.cache.IndexTask(context.Background(), task); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) TaskDeleted(id string) error {
	if err := r.cache.DeleteTask(context.Background(), id); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) TaskUpdated(task internal.Tasks) error {
	if err := r.cache.UpdateTask(context.Background(), task); err != nil {
		return err
	}
	return nil
}
