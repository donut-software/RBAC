package events

import (
	"context"
	"encoding/json"
	"rbac/internal"
	"strings"

	"github.com/go-redis/redis/v8"
)

func (r *RBACEvents) AccountCreated(msg *redis.Message) error {
	var account internal.Account
	if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&account); err != nil {
		return err
	}
	if err := r.cache.IndexAccount(context.Background(), account); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) AccountDeleted(msg *redis.Message) error {
	var id string
	if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&id); err != nil {
		return err
	}
	if err := r.cache.DeleteAccount(context.Background(), id); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) AccountUpdated(msg *redis.Message) error {
	var profile internal.Profile
	if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&profile); err != nil {
		return err
	}
	if err := r.cache.UpdateProfile(context.Background(), profile); err != nil {
		return err
	}
	return nil
}
