package worker

import (
	"context"
	"log"
	"sync"

	"github.com/c4erries/Sentry/internal/model"
)

type EventHandler interface {
	Process(ctx context.Context, e model.Event) error
}

func StartPool(ctx context.Context, jobs <-chan model.Event, handler EventHandler, workerCount int) *sync.WaitGroup {
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
					if err := handler.Process(ctx, job); err != nil {
						log.Printf("[worker-%d] error processing event %v: %v", workerID, job, err)
					}
				}
			}
		}(id)
	}

	return &wg
}
