package redis

import (
	"context"
	"rbac/internal"
)

// Created publishes a message indicating a accounts was created.
func (t *RBAC) AccountCreated(ctx context.Context, accounts internal.Account) error {
	return t.publish(ctx, "Account.Created", internal.EVENT_ACCOUNT_CREATED, accounts)
}

// Deleted publishes a message indicating a accounts was deleted.
func (t *RBAC) AccountDeleted(ctx context.Context, id string) error {
	return t.publish(ctx, "Account.Deleted", internal.EVENT_ACCOUNT_DELETED, id)
}

// Updated publishes a message indicating a accounts was updated.
func (t *RBAC) AccountUpdated(ctx context.Context, profile internal.Account) error {
	return t.publish(ctx, "Account.Updated", internal.EVENT_ACCOUNT_UPDATED, profile)
}
