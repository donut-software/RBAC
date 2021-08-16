package redis

import (
	"context"
	"rbac/internal"
)

// Created publishes a message indicating a profiles was created.
func (t *RBAC) ProfileCreated(ctx context.Context, profile internal.Profile) error {
	return t.publish(ctx, "Profile.Created", "profiles.event.created", profile)
}

// Deleted publishes a message indicating a profiles was deleted.
func (t *RBAC) ProfileDeleted(ctx context.Context, id string) error {
	return t.publish(ctx, "Profile.Deleted", "profiles.event.deleted", id)
}

// Updated publishes a message indicating a profiles was updated.
func (t *RBAC) ProfileUpdated(ctx context.Context, profile internal.Profile) error {
	return t.publish(ctx, "Profile.Updated", "profiles.event.updated", profile)
}
