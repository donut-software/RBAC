package kafka

import (
	"context"
	"rbac/internal"
)

// Created publishes a message indicating a accountRole was created.
func (t *RBAC) AccountRoleCreated(ctx context.Context, accountRole internal.AccountRoles) error {
	return t.publish(ctx, "AccountRole.Created", internal.EVENT_ACCOUNTROLE_CREATED, accountRole)
}

// Deleted publishes a message indicating a accountRole was deleted.
func (t *RBAC) AccountRoleDeleted(ctx context.Context, id string) error {
	return t.publish(ctx, "AccountRole.Deleted", internal.EVENT_ACCOUNTROLE_DELETED, id)
}

// Updated publishes a message indicating a accountRole was updated.
func (t *RBAC) AccountRoleUpdated(ctx context.Context, roleTask internal.AccountRoles) error {
	return t.publish(ctx, "AccountRole.Updated", internal.EVENT_ACCOUNTROLE_UPDATED, roleTask)
}
