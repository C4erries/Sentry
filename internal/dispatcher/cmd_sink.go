package dispatcher

import (
	"context"
	"log"

	"github.com/c4erries/Sentry/internal/model"
)

type CmdSink struct {
}

func NewCmdSink() *CmdSink {
	return &CmdSink{}
}

type Printable interface {
	String() string
	IsPrintable() bool
}

func (s *CmdSink) SendAlert(ctx context.Context, alert *model.Alert) error {
	events_string := ""
	for _, event := range alert.Events {
		events_string += event.BaseEvent.EventType.String() + " "
	}
	log.Printf("ALERT | Rule: %s, Level: %s, DetectedAt: %s, Events: %s", alert.Rule, alert.Level, alert.DetectedAt.UTC().String(), events_string)
	data, ok := alert.Data.(Printable)
	if !ok {
		return nil
	}
	if data.IsPrintable() {
		log.Printf("####Data####")
		log.Printf("%s", data.String())
		log.Printf("############")
	}
	return nil
}

func (s *CmdSink) ID() string {
	return "cmd_sink"
}
