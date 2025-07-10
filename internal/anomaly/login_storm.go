package anomaly

import (
	"context"
	"fmt"
	"time"

	"github.com/c4erries/Sentry/internal/model"
	"github.com/c4erries/Sentry/internal/redis"
)

type LoginStormDetector struct {
	redis     redis.RedisClient
	window    time.Duration
	threshold int64
	prefix    string
}

func NewLoginStormDetector(r redis.RedisClient, window time.Duration, threshold int64) *LoginStormDetector {
	return &LoginStormDetector{
		redis:     r,
		window:    window,
		threshold: threshold,
		prefix:    "login",
	}
}

func (d *LoginStormDetector) ID() string {
	return "login_storm"
}

func (d *LoginStormDetector) Process(ctx context.Context, e *model.Event) (*model.Alert, error) {
	if e.EventType != model.EventLogin {
		return nil, nil
	}

	key := fmt.Sprintf("%s:%s", d.prefix, e.UserId)
	nowTS := e.Timestamp.Unix()

	d.redis.ZAdd(ctx, key, redis.Z{
		Score:  float64(nowTS),
		Member: e.ID,
	})

	minScore := float64(nowTS) - d.window.Seconds()
	d.redis.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%f", minScore))

	d.redis.Expire(ctx, key, d.window+redisTTLBuffer)

	count, err := d.redis.ZCard(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis ZRANGEBYSCORE error: %v", err)
	}

	if count <= d.threshold {
		return nil, nil
	}
	eventIDs, err := d.redis.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: fmt.Sprintf("%f", minScore),
		Max: fmt.Sprintf("%f", float64(nowTS)),
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("redis ZRANGEBYSCORE error: %v", err)
	}

	return model.NewAlert(
		model.AnomalyLoginStorm,
		[]*model.Event{e},
		model.AlertWarning,
		time.Now(),
		model.LoginStormData{
			EventIDs: eventIDs,
		},
	), nil

}
