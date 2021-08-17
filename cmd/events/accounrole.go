package events

import (
	"context"
	"encoding/json"
	"rbac/internal"
	"strings"

	"github.com/go-redis/redis/v8"
)

func (r *RBACEvents) AccountRoleCreated(msg *redis.Message) error {
	var accountRole internal.AccountRoles
	if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&accountRole); err != nil {
		return err
	}
	if err := r.cache.IndexAccountRole(context.Background(), accountRole); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) AccountRoleDeleted(msg *redis.Message) error {
	var id string
	if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&id); err != nil {
		return err
	}
	if err := r.cache.DeleteAccountRole(context.Background(), id); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) AccountRoleUpdated(msg *redis.Message) error {
	var accountRole internal.AccountRoles
	if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&accountRole); err != nil {
		return err
	}
	if err := r.cache.UpdateAccountRole(context.Background(), accountRole); err != nil {
		return err
	}
	return nil
}
