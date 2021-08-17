package rabbitmq

import (
	"context"
	"rbac/internal"
)

// Created publishes a message indicating a roleTasks was created.
func (t *RBAC) RoleTaskCreated(ctx context.Context, roleTasks internal.RoleTasks) error {
	return t.publish(ctx, "RoleTask.Created", internal.EVENT_ROLETASK_CREATED, roleTasks)
}

// Deleted publishes a message indicating a roleTasks was deleted.
func (t *RBAC) RoleTaskDeleted(ctx context.Context, id string) error {
	return t.publish(ctx, "RoleTask.Deleted", internal.EVENT_ROLETASK_DELETED, id)
}

// Updated publishes a message indicating a roleTasks was updated.
func (t *RBAC) RoleTaskUpdated(ctx context.Context, roleTask internal.RoleTasks) error {
	return t.publish(ctx, "RoleTask.Updated", internal.EVENT_ROLETASK_UPDATED, roleTask)
}
