package dispatcher

import (
	"context"
	"log/slog"

	"github.com/c4erries/Sentry/internal/model"
)

type CmdSink struct {
}

func NewCmdSink() *CmdSink {
	return &CmdSink{}
}

func (s *CmdSink) SendAlert(ctx context.Context, alert *model.Alert) error {
	slog.InfoContext(ctx, "Alert", slog.Any("alert", alert))
	return nil
}

func (s *CmdSink) ID() string {
	return "cmd_sink"
}
