package kafka

import (
	"context"
	"encoding/json"
	"log"
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
				log.Println("Context canceled, stoping consumer")
				break
			}
			log.Printf("Fetch error: %v", err)
			continue
		}

		e := new(model.Event)
		if err = json.Unmarshal(m.Value, e); err != nil {
			log.Printf("json unmarshal event error: %v", err)
			continue
		}
		if err = e.Normalize(); err != nil {
			log.Printf("[event-%s] normalization error: %v", e.ID, err)
		}
		if err = e.Validate(); err != nil {
			log.Printf("[event-%s] validation failed: %v", e.ID, err)
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
