package dispatcher

import (
	"context"
	"log"
	"sync"

	"github.com/c4erries/Sentry/internal/model"
)

type EventDispatcher struct {
	Sinks []EventSink
	Chan  chan *model.Event
}

func NewEventDispatcher(sinks []EventSink, bufferSize int) *EventDispatcher {
	return &EventDispatcher{
		Sinks: sinks,
		Chan:  make(chan *model.Event, bufferSize),
	}
}

func (d *EventDispatcher) SendAll(ctx context.Context, event *model.Event) {
	var wg sync.WaitGroup
	wg.Add(len(d.Sinks))

	for _, sink := range d.Sinks {
		go func(s EventSink) {
			defer wg.Done()
			if err := s.SendEvent(ctx, event); err != nil {
				log.Printf("[event-dispatcher] sink-%s failed: %v", s.ID(), err)
			}
		}(sink)
	}

	wg.Wait()
}

func (d *EventDispatcher) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("[event-dispatcher] shutdown requested")
			return

		case event, ok := <-d.Chan:
			if !ok {
				log.Println("[event-dispatcher] event channel closed")
				return
			}
			d.SendAll(ctx, event)
		}
	}
}
