package kafka

import (
	"context"
	"rbac/internal"
)

// Created publishes a message indicating a roles was created.
func (t *RBAC) RoleCreated(ctx context.Context, roles internal.Roles) error {
	return t.publish(ctx, "Role.Created", internal.EVENT_ROLE_CREATED, roles)
}

// Deleted publishes a message indicating a roles was deleted.
func (t *RBAC) RoleDeleted(ctx context.Context, id string) error {
	return t.publish(ctx, "Role.Deleted", internal.EVENT_ROLETASK_DELETED, id)
}

// Updated publishes a message indicating a roles was updated.
func (t *RBAC) RoleUpdated(ctx context.Context, role internal.Roles) error {
	return t.publish(ctx, "Role.Updated", internal.EVENT_ROLETASK_UPDATED, role)
}
