package redis

import (
	"context"
	"rbac/internal"
)

// Created publishes a message indicating a roles was created.
func (t *RBAC) RoleCreated(ctx context.Context, roles internal.Roles) error {
	return t.publish(ctx, "Role.Created", "roles.event.created", roles)
}

// Deleted publishes a message indicating a roles was deleted.
func (t *RBAC) RoleDeleted(ctx context.Context, id string) error {
	return t.publish(ctx, "Role.Deleted", "roles.event.deleted", id)
}

// Updated publishes a message indicating a roles was updated.
func (t *RBAC) RoleUpdated(ctx context.Context, role internal.Roles) error {
	return t.publish(ctx, "Role.Updated", "roles.event.updated", role)
}
