package alert

import (
	"context"
	"log"
	"sync"

	"github.com/c4erries/Sentry/internal/model"
)

type Dispatcher struct {
	Sinks []AlertSink
	Chan  chan *model.Alert
}

func NewDispatcher(sinks []AlertSink, bufferSize int) *Dispatcher {
	return &Dispatcher{
		Sinks: sinks,
		Chan:  make(chan *model.Alert, bufferSize),
	}
}

func (d *Dispatcher) SendAll(ctx context.Context, alert *model.Alert) {
	var wg sync.WaitGroup
	wg.Add(len(d.Sinks))

	for _, sink := range d.Sinks {
		go func(s AlertSink) {
			defer wg.Done()
			if err := s.Send(ctx, alert); err != nil {
				log.Printf("[dispatcher] sink-%s failed: %v", s.ID(), err)
			}
		}(sink)
	}

	wg.Wait()
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
			d.SendAll(ctx, alert)
		}
	}
}
