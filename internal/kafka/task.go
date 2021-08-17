package kafka

import (
	"context"
	"rbac/internal"
)

// Created publishes a message indicating a tasks was created.
func (t *RBAC) TaskCreated(ctx context.Context, tasks internal.Tasks) error {
	return t.publish(ctx, "Task.Created", internal.EVENT_TASK_CREATED, tasks)
}

// Deleted publishes a message indicating a tasks was deleted.
func (t *RBAC) TaskDeleted(ctx context.Context, id string) error {
	return t.publish(ctx, "Task.Deleted", internal.EVENT_TASK_UPDATED, id)
}

// Updated publishes a message indicating a tasks was updated.
func (t *RBAC) TaskUpdated(ctx context.Context, task internal.Tasks) error {
	return t.publish(ctx, "Task.Updated", internal.EVENT_TASK_UPDATED, task)
}
