package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/c4erries/Sentry/internal/model"
	"github.com/segmentio/kafka-go"
)

const (
	sessionTimeout = 10 * time.Second
)

type Message interface {
	kafka.Message
}

type KafkaEvent struct {
	*model.Event
	Commit func() error
}

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic string, groupID string) (*Consumer, error) {
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
	return &Consumer{reader: r}, nil
}

func (c *Consumer) Start(ctx context.Context, out chan *KafkaEvent) {
	defer c.reader.Close()

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				slog.InfoContext(ctx, "Context canceled, stoping consumer")
				break
			}
			slog.WarnContext(ctx, "Fetch kafka message error", slog.Any("err", err))
			continue
		}

		e := new(model.Event)
		if err = json.Unmarshal(m.Value, e); err != nil {
			slog.ErrorContext(ctx, "json unmarshal event error", slog.Any("err", err))
			continue
		}
		if err = e.Normalize(); err != nil {
			slog.ErrorContext(ctx, "normalization error", slog.String("event_id", e.ID), slog.Any("err", err))
		}
		if err = e.Validate(); err != nil {
			slog.ErrorContext(ctx, "validation failed", slog.String("event_id", e.ID), slog.Any("err", err))
			continue
		}

		kevent := &KafkaEvent{
			Event: e,
			Commit: func() error {
				return c.reader.CommitMessages(ctx, m)
			},
		}

		out <- kevent
	}
}
