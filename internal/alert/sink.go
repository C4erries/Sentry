package alert

import (
	"context"

	"github.com/c4erries/Sentry/internal/model"
)

type AlertSink interface {
	Send(ctx context.Context, alert model.Alert) error
}
