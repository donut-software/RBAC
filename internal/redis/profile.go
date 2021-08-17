package redis

import (
	"context"
	"rbac/internal"
)

// Created publishes a message indicating a profiles was created.
func (t *RBAC) ProfileCreated(ctx context.Context, profile internal.Profile) error {
	return t.publish(ctx, "Profile.Created", internal.EVENT_PROFILE_CREATED, profile)
}

// Deleted publishes a message indicating a profiles was deleted.
func (t *RBAC) ProfileDeleted(ctx context.Context, id string) error {
	return t.publish(ctx, "Profile.Deleted", internal.EVENT_PROFILE_DELETED, id)
}

// Updated publishes a message indicating a profiles was updated.
func (t *RBAC) ProfileUpdated(ctx context.Context, profile internal.Profile) error {
	return t.publish(ctx, "Profile.Updated", internal.EVENT_PROFILE_UPDATED, profile)
}
