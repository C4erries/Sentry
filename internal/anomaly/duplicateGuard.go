package anomaly

import (
	"context"
	"fmt"
	"time"

	"github.com/c4erries/Sentry/internal/model"
	"github.com/c4erries/Sentry/internal/redis"
)

type DuplicateGuard struct {
	redis  redis.RedisClient
	ttl    time.Duration
	prefix string
}

func NewDuplicateGuard(r redis.RedisClient, ttl time.Duration) *DuplicateGuard {
	return &DuplicateGuard{
		redis:  r,
		ttl:    ttl,
		prefix: "event",
	}
}

func (d *DuplicateGuard) IsDuplicate(ctx context.Context, e model.Event) (bool, error) {
	key := fmt.Sprintf("%s:%s", d.prefix, e.ID)

	ok, err := d.redis.SetNX(ctx, key, "seen", d.ttl).Result()
	if err != nil {
		return false, fmt.Errorf("redis SETNX failed: %v", err)
	}
	return !ok, nil
}
