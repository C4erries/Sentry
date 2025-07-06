package alert

import (
	"context"
	"log"

	"github.com/c4erries/Sentry/internal/model"
)

type Dispatcher struct {
	Sink AlertSink
	Chan chan model.Alert
}

func NewDispatcher(sink AlertSink, bufferSize int) *Dispatcher {
	return &Dispatcher{
		Sink: sink,
		Chan: make(chan model.Alert, bufferSize),
	}
}

func (d *Dispatcher) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("[dispatcher] shutdown requested")
			return

		case alert, ok := <-d.Chan:
			if !ok {
				log.Println("[dispatcher] alert channel closed")
				return
			}
			if err := d.Sink.Send(ctx, alert); err != nil {
				log.Printf("[dispatcher] failed to send alert: %v", err)
			}
		}
	}
}
