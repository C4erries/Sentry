package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/c4erries/Sentry/internal/events"
	"github.com/segmentio/kafka-go"
)

const (
	sessionTimeout = 10 * time.Second
)

type Message interface {
	kafka.Message
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
	})
	return &Consumer{reader: r}, nil
}

func (c *Consumer) Start(ctx context.Context, out chan events.Event) {
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

		var e events.Event
		if err = json.Unmarshal(m.Value, &e); err != nil {
			log.Printf("json unmarshal event error: %v", err)
			continue
		}

		out <- e
		if err := c.reader.CommitMessages(ctx, m); err != nil {
			log.Printf("commit error: %v", err)
		}
	}
}
