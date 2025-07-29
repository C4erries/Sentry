package kafka

import (
	"context"
	"encoding/json"
	"errors"
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

//go:generate go run github.com/vektra/mockery/v2@v2.53.4 --name=Reader
type Reader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type Consumer struct {
	reader Reader
}

func NewConsumer(reader Reader) (*Consumer, error) {
	return &Consumer{reader: reader}, nil
}

func (c *Consumer) Start(ctx context.Context, out chan *KafkaEvent) {
	defer c.reader.Close()

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			// Если контекст отменён — считаем это нормальным завершением
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				slog.InfoContext(ctx, "Context canceled, stopping consumer")
				break
			}
			// Для других ошибок логируем предупреждение и продолжаем попытки чтения
			slog.WarnContext(ctx, "Fetch kafka message error", slog.Any("err", err))
			continue
		}

		e := new(model.Event)
		if err = json.Unmarshal(m.Value, e); err != nil {
			slog.ErrorContext(ctx, "json unmarshal event error", slog.Any("err", err))
			continue
		}
		if err = e.Normalize(); err != nil {
			slog.ErrorContext(ctx, "normalization error", slog.String("event_id", e.ID.String()), slog.Any("err", err))
		}
		if err = e.Validate(); err != nil {
			slog.ErrorContext(ctx, "validation failed", slog.String("event_id", e.ID.String()), slog.Any("err", err))
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
