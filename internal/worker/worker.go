package worker

import (
	"context"
	"log/slog"
	"sync"

	"github.com/c4erries/Sentry/internal/kafka"
	"github.com/c4erries/Sentry/internal/model"
)

type EventHandler interface {
	Process(ctx context.Context, e *model.Event) error
}

func StartPool(ctx context.Context, jobs <-chan *kafka.KafkaEvent, handler EventHandler, workerCount int) *sync.WaitGroup {
	var wg sync.WaitGroup

	for id := 0; id < workerCount; id++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			slog.InfoContext(ctx, "worker started", slog.Int("workerID", workerID))
			for {
				select {
				case <-ctx.Done():
					slog.InfoContext(ctx, "worker stopping (context canceled)", slog.Int("workerID", workerID))
					return

				case job, ok := <-jobs:
					if !ok {
						slog.InfoContext(ctx, "worker stopping (jobs channel closed)", slog.Int("workerID", workerID))
						return
					}
					if err := handler.Process(ctx, job.Event); err != nil {
						slog.InfoContext(ctx, "worker error processing event",
							slog.Int("workerID", workerID),
							slog.Any("event", job.Event),
							slog.Any("err", err),
						)
					} else if err := job.Commit(); err != nil {
						slog.InfoContext(ctx, "event commit error",
							slog.Any("job", job),
							slog.Any("err", err),
						)
					}
				}
			}
		}(id)
	}

	return &wg
}
