package events

import (
	"context"
	"encoding/json"
	"rbac/internal"
	"strings"

	"github.com/go-redis/redis/v8"
)

func (r *RBACEvents) RoleCreated(msg *redis.Message) error {
	var role internal.Roles
	if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&role); err != nil {
		return err
	}
	if err := r.cache.IndexRole(context.Background(), role); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) RoleDeleted(msg *redis.Message) error {
	var id string
	if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&id); err != nil {
		return err
	}
	if err := r.cache.DeleteRole(context.Background(), id); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) RoleUpdated(msg *redis.Message) error {
	var role internal.Roles
	if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&role); err != nil {
		return err
	}
	if err := r.cache.UpdateRole(context.Background(), role); err != nil {
		return err
	}
	return nil
}
