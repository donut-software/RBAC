package events

import (
	"context"
	"encoding/json"
	"rbac/internal"
	"strings"

	"github.com/go-redis/redis/v8"
)

func (r *RBACEvents) RoleTaskCreated(msg *redis.Message) error {
	var roleTask internal.RoleTasks
	if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&roleTask); err != nil {
		return err
	}
	if err := r.cache.IndexRoleTask(context.Background(), roleTask); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) RoleTaskDeleted(msg *redis.Message) error {
	var id string
	if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&id); err != nil {
		return err
	}
	if err := r.cache.DeleteRoleTask(context.Background(), id); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) RoleTaskUpdated(msg *redis.Message) error {
	var roleTask internal.RoleTasks
	if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&roleTask); err != nil {
		return err
	}
	if err := r.cache.UpdateRoleTask(context.Background(), roleTask); err != nil {
		return err
	}
	return nil
}
