package dispatcher

import (
	"context"
	"log"
	"sync"

	"github.com/c4erries/Sentry/internal/model"
)

type AlertDispatcher struct {
	Sinks []AlertSink
	Chan  chan *model.Alert
}

func NewAlertDispatcher(sinks []AlertSink, bufferSize int) *AlertDispatcher {
	return &AlertDispatcher{
		Sinks: sinks,
		Chan:  make(chan *model.Alert, bufferSize),
	}
}

func (d *AlertDispatcher) SendAll(ctx context.Context, alert *model.Alert) {
	var wg sync.WaitGroup
	wg.Add(len(d.Sinks))

	for _, sink := range d.Sinks {
		go func(s AlertSink) {
			defer wg.Done()
			if err := s.SendAlert(ctx, alert); err != nil {
				log.Printf("[alert-dispatcher] sink-%s failed: %v", s.ID(), err)
			}
		}(sink)
	}

	wg.Wait()
}

func (d *AlertDispatcher) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("[alert-dispatcher] shutdown requested")
			return

		case alert, ok := <-d.Chan:
			if !ok {
				log.Println("[alert-dispatcher] alert channel closed")
				return
			}
			d.SendAll(ctx, alert)
		}
	}
}
