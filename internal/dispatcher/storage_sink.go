package dispatcher

import (
	"context"

	"github.com/c4erries/Sentry/internal/model"
)

type EventRepository interface {
	Save(ctx context.Context, event *model.Event) error
}

type AlertRepository interface {
	Save(ctx context.Context, alert *model.Alert) error
}

type StorageSink struct {
	eventRepo EventRepository
	alertRepo AlertRepository
}

func NewStorageSink(eventRepo EventRepository, alertRepo AlertRepository) *StorageSink {
	return &StorageSink{
		eventRepo: eventRepo,
		alertRepo: alertRepo,
	}
}

func (s *StorageSink) ID() string {
	return "storage_sink"
}

func (s *StorageSink) SendEvent(ctx context.Context, event *model.Event) error {
	return s.eventRepo.Save(ctx, event)
}

func (s *StorageSink) SendAlert(ctx context.Context, alert *model.Alert) error {
	return s.alertRepo.Save(ctx, alert)
}
