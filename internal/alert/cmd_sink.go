package alert

import (
	"context"
	"log"

	"github.com/c4erries/Sentry/internal/model"
)

type CmdSink struct {
}

func (s *CmdSink) Send(ctx context.Context, alert model.Alert) error {
	events_string := ""
	for _, event := range alert.Events {
		events_string += event.BaseEvent.EventType.String() + " "
	}
	log.Printf("ALERT | Rule: %s, Level: %s, DetectedAt: %s, Events: %s", alert.Rule, alert.Level, alert.DetectedAt.UTC().String(), events_string)
	return nil
}
