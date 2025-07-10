package worker

import (
	"context"
	"log"
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
			log.Printf("[worker-%d] started", workerID)
			for {
				select {
				case <-ctx.Done():
					log.Printf("[worker-%d] stopping (context canceled)", workerID)
					return

				case job, ok := <-jobs:
					if !ok {
						log.Printf("[worker-%d] stopping (jobs channel closed)", workerID)
						return
					}
					if err := handler.Process(ctx, job.Event); err != nil {
						log.Printf("[worker-%d] error processing event %v: %v", workerID, job, err)
					}
					if err := job.Commit(); err != nil {
						log.Printf("[event-%s] commit error: %v", job.ID, err)
					}
				}
			}
		}(id)
	}

	return &wg
}
