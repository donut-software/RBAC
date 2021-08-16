package redis

import (
	"bytes"
	"context"
	"encoding/json"
	"rbac/internal"

	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

// Account represents the repository used for publishing Account records.
type RBAC struct {
	client *redis.Client
}

// NewAccount instantiates the Account repository.
func NewAccount(client *redis.Client) *RBAC {
	return &RBAC{
		client: client,
	}
}
func (t *RBAC) publish(ctx context.Context, spanName, channel string, e interface{}) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, spanName)
	defer span.End()

	span.SetAttributes(
		semconv.DBSystemRedis,
		attribute.KeyValue{
			Key:   "db.statement",
			Value: attribute.StringValue("PUBLISH"),
		},
	)

	//-

	var b bytes.Buffer

	if err := json.NewEncoder(&b).Encode(e); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.Encode")
	}

	res := t.client.Publish(ctx, channel, b.Bytes())
	if err := res.Err(); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Publish")
	}

	return nil
}
