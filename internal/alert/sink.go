package alert

import (
	"context"

	"github.com/c4erries/Sentry/internal/model"
)

type AlertSink interface {
	ID() string
	Send(ctx context.Context, alert *model.Alert) error
}
