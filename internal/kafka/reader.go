package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaReader struct {
	inner *kafka.Reader
}

func NewReader(brokers []string, topic string, groupID string) (*KafkaReader, error) {
	if groupID == "" {
		return nil, fmt.Errorf("group id can't be empty")
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		SessionTimeout: sessionTimeout,
		CommitInterval: 0,
	})
	return &KafkaReader{inner: r}, nil
}

func (r *KafkaReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return r.inner.ReadMessage(ctx)
}

func (r *KafkaReader) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	return r.inner.CommitMessages(ctx, msgs...)
}

func (r *KafkaReader) Close() error {
	return r.inner.Close()
}
