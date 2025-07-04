package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/c4erries/Sentry/internal/events"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(address []string, topic string) (*Producer, error) {
	w := &kafka.Writer{
		Addr:         kafka.TCP(address...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireAll,
	}
	return &Producer{writer: w}, nil
}

func (p *Producer) Produce(ctx context.Context, e events.Event) error {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("event marshal error: %v", err)
	}

	msg := kafka.Message{
		Key:   []byte(e.UserId),
		Value: data,
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) ProduceBatch(ctx context.Context, evs ...events.Event) error {
	var msgs []kafka.Message
	for _, e := range evs {
		data, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("event marshal error: %v", err)
		}

		msg := kafka.Message{
			Key:   []byte(e.UserId),
			Value: data,
		}
		msgs = append(msgs, msg)
	}
	if err := p.writer.WriteMessages(ctx, msgs...); err != nil {
		return fmt.Errorf("write message error: %v", err)
	}
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
