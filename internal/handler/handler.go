package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/c4erries/Sentry/internal/events"
)

type Processor struct {
	// db *sql.DB
	// cache redis.Client
}

func NewProcessor() (*Processor, error) {
	return &Processor{}, nil
}

func (p *Processor) Process(ctx context.Context, e events.Event) error {
	log.Printf("Processing [type:%v] [timestamp:%v]", e.EventType, e.Timestamp.String())
	switch e.EventType {
	case events.EventLogin:
		return p.handleLogin(ctx, e)
	case events.EventTransaction:
		return p.handleTransaction(ctx, e)
	default:
		return fmt.Errorf("can't process event: Unimplemented event type")
	}
}
