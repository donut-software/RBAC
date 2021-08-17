package events

import (
	"context"
	"encoding/json"
	"rbac/internal"
	"strings"

	"github.com/go-redis/redis/v8"
)

func (r *RBACEvents) TaskCreated(msg *redis.Message) error {
	var task internal.Tasks
	if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&task); err != nil {
		return err
	}
	if err := r.cache.IndexTask(context.Background(), task); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) TaskDeleted(msg *redis.Message) error {
	var id string
	if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&id); err != nil {
		return err
	}
	if err := r.cache.DeleteTask(context.Background(), id); err != nil {
		return err
	}
	return nil
}
func (r *RBACEvents) TaskUpdated(msg *redis.Message) error {
	var task internal.Tasks
	if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&task); err != nil {
		return err
	}
	if err := r.cache.UpdateTask(context.Background(), task); err != nil {
		return err
	}
	return nil
}
