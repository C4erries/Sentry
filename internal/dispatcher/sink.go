package dispatcher

import (
	"context"

	"github.com/c4erries/Sentry/internal/model"
)

type EventSink interface {
	ID() string
	SendEvent(ctx context.Context, event *model.Event) error
}

type AlertSink interface {
	ID() string
	SendAlert(ctx context.Context, alert *model.Alert) error
}
