package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"rbac/internal"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

// Account represents the repository used for publishing Account records.
type RBAC struct {
	producer  *kafka.Producer
	topicName string
}

type event struct {
	Type  string
	Value interface{}
}

// NewAccount instantiates the Account repository.
func NewRBAC(producer *kafka.Producer, topicName string) *RBAC {
	return &RBAC{
		topicName: topicName,
		producer:  producer,
	}
}

func (t *RBAC) publish(ctx context.Context, spanName, msgType string, e interface{}) error {
	_, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, spanName)
	defer span.End()

	span.SetAttributes(
		attribute.KeyValue{
			Key:   semconv.MessagingSystemKey,
			Value: attribute.StringValue("kafka"),
		},
	)

	//-

	var b bytes.Buffer

	evt := event{
		Type:  msgType,
		Value: e,
	}

	if err := json.NewEncoder(&b).Encode(evt); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.Encode")
	}

	if err := t.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &t.topicName,
			Partition: kafka.PartitionAny,
		},
		Value: b.Bytes(),
	}, nil); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "product.Producer")
	}

	return nil
}
