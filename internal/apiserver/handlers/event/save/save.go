package save

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/c4erries/Sentry/internal/model"
)

type EventSaver interface {
	Save(ctx context.Context, event *model.Event) error
}

func New(log *slog.Logger, repo EventSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
