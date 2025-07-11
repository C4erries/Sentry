package dispatcher

import (
	"context"
	"fmt"

	"github.com/c4erries/Sentry/internal/model"
)

type AlertCache interface {
	SaveAlert(ctx context.Context, alert *model.Alert) error
}

type Publisher interface {
	Publish(ctx context.Context, alert *model.Alert) error
}

type RedisSink struct {
	cache     AlertCache
	publisher Publisher
}

func NewRedisSink(cache AlertCache, publisher Publisher) *RedisSink {
	return &RedisSink{
		cache:     cache,
		publisher: publisher,
	}
}

func (s *RedisSink) ID() string {
	return "redis_sink"
}

func (s *RedisSink) SendAlert(ctx context.Context, alert *model.Alert) error {
	err := s.cache.SaveAlert(ctx, alert)
	if err != nil {
		return fmt.Errorf("RedisSink can't SEND [alert-%s]: %v", alert.ID, err)
	}
	err = s.publisher.Publish(ctx, alert)
	if err != nil {
		return fmt.Errorf("RedisSink can't PUBLISH [alert-%s]: %v", alert.ID, err)
	}
	return nil
}
