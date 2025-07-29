package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/c4erries/Sentry/internal/model"
	"github.com/segmentio/kafka-go"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.4 --name=Writer
type Writer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type Producer struct {
	writer Writer
}

func NewProducer(writer Writer) (*Producer, error) {
	return &Producer{writer: writer}, nil
}

func (p *Producer) Produce(ctx context.Context, e *model.Event) error {
	if err := e.Validate(); err != nil {
		return fmt.Errorf("event:%s -- validation failed:%v", e.ID, err)
	}
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("event marshal error: %v", err)
	}

	msg := kafka.Message{
		Key:   []byte(e.UserID),
		Value: data,
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) ProduceBatch(ctx context.Context, evs ...*model.Event) error {
	var msgs []kafka.Message
	for _, e := range evs {

		if err := e.Validate(); err != nil {
			return fmt.Errorf("event:%s -- validation failed:%v", e.ID, err)
		}

		data, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("event marshal error: %v", err)
		}

		msg := kafka.Message{
			Key:   []byte(e.UserID),
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
