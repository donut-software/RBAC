package rabbitmq

import (
	"bytes"
	"context"
	"encoding/gob"
	"rbac/internal"
	"time"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

// Account represents the repository used for publishing Account records.
type RBAC struct {
	ch *amqp.Channel
}

// NewTask instantiates the Account repository.
func NewRBAC(channel *amqp.Channel) (*RBAC, error) {
	return &RBAC{
		ch: channel,
	}, nil
}

func (t *RBAC) publish(ctx context.Context, spanName, routingKey string, e interface{}) error {
	_, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, spanName)
	defer span.End()
	span.SetAttributes(
		attribute.KeyValue{
			Key:   semconv.MessagingSystemKey,
			Value: attribute.StringValue("rabbitmq"),
		},
		attribute.KeyValue{
			Key:   semconv.MessagingRabbitMQRoutingKeyKey,
			Value: attribute.StringValue(routingKey),
		},
	)

	//-

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(e); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.Encode")
	}

	err := t.ch.Publish(
		"rbac",     // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			AppId:       "rbac-rest-server",
			ContentType: "application/x-encoding-gob", // XXX: We will revisit this in future episodes
			Body:        b.Bytes(),
			Timestamp:   time.Now(),
		})
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "ch.Publish")
	}

	return nil
}
