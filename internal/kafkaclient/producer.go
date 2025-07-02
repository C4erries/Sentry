package kafkaclient

import (
	"context"
	"time"

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

func (p *Producer) Produce(ctx context.Context, msg ...kafka.Message) error {
	return p.writer.WriteMessages(ctx, msg...)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
