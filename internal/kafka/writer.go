package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	inner *kafka.Writer
}

func NewWriter(address []string, topic string) (*KafkaWriter, error) {
	w := &kafka.Writer{
		Addr:         kafka.TCP(address...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireAll,
	}
	return &KafkaWriter{inner: w}, nil
}

func (w *KafkaWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	return w.inner.WriteMessages(ctx, msgs...)
}

func (w *KafkaWriter) Close() error {
	return w.inner.Close()
}
