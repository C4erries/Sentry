package dispatcher

import (
	"context"
	"sync"
	"time"

	"github.com/c4erries/Sentry/internal/model"
)

type DeferredAlert struct {
	Alert      *model.Alert
	RetryAt    time.Time
	RetryCount int
}

type Repository interface {
	AllEventsExist(ctx context.Context, eventIDs []string) (bool, error)
	Publish(ctx context.Context)
}

type AlertRetryQueue struct {
	mu         sync.Mutex
	queue      []*DeferredAlert
	repository Repository
}

func NewAlertRetryQueue(repository Repository) *AlertRetryQueue {
	return &AlertRetryQueue{
		mu:         sync.Mutex{},
		queue:      make([]*DeferredAlert, 0),
		repository: repository,
	}
}

func (q *AlertRetryQueue) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			q.retryDueAlerts(ctx)
		}
	}
}

func (q *AlertRetryQueue) retryDueAlerts(ctx context.Context) {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := time.Now()
	newQueue := []*DeferredAlert{}

	for _, d := range q.queue {
		if d.RetryAt.After(now) {
			newQueue = append(newQueue, d)
			continue
		}
		panic("unimplemented")
		/*
			exists, err := q.repository.AllEventsExist(ctx, d.Alert.Event)
			if exists {

			} else {
				d.RetryCount++
				if d.RetryCount < 5 {
					d.RetryAt = now.Add(time.Duration(d.RetryCount*500) * time.Millisecond)
					newQueue = append(newQueue, d)
				} else {
					slog.Error("Max retries exceeded", slog.Any("alert", d.Alert))
				}
			}
		*/
	}

	q.queue = newQueue
}
