package anomaly

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/c4erries/Sentry/internal/model"
	"github.com/c4erries/Sentry/internal/redis"
)

type GeoCache interface {
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
}

type GeoSwitchingDetector struct {
	redis  GeoCache
	window time.Duration
	prefix string
}

func NewGeoSwitchingDetector(r GeoCache, window time.Duration) *GeoSwitchingDetector {
	return &GeoSwitchingDetector{
		redis:  r,
		window: window,
		prefix: "geo",
	}
}

func (d *GeoSwitchingDetector) ID() string {
	return "geo_swithcing"
}

func (d *GeoSwitchingDetector) Process(ctx context.Context, e *model.Event) (*model.Alert, error) {
	currCountry := e.GeoCountry // ?
	nowTS := e.Timestamp.Unix()

	key := fmt.Sprintf("%s:%s", d.prefix, e.UserId)

	exists, err := d.redis.Exists(ctx, key).Result()
	if exists == 0 {
		d.redis.HSet(ctx, key, map[string]interface{}{
			"country":   currCountry,
			"timestamp": nowTS,
			"event_id":  e.ID,
		})
		d.redis.Expire(ctx, key, d.window+redisTTLBuffer)
		return nil, nil
	}

	prev, err := d.redis.HMGet(ctx, key, "country", "timestamp", "event_id").Result()
	if err != nil {
		return nil, fmt.Errorf("redis HMGet error: %v", err)
	}

	prevCountry, ok := prev[0].(string)
	if !ok {
		return nil, fmt.Errorf("[redis] can't find field 'country' on [event-%s]", e.ID)
	}
	prevTSStr, ok := prev[1].(string)
	if !ok {
		return nil, fmt.Errorf("[redis] can't find field 'timestamp' on [event-%s]", e.ID)
	}
	prevID, ok := prev[2].(string)
	if !ok {
		return nil, fmt.Errorf("[redis] can't find field 'event_id' on [event-%s]", e.ID)
	}

	d.redis.HSet(ctx, key, map[string]interface{}{
		"country":   currCountry,
		"timestamp": nowTS,
		"event_id":  e.ID,
	})

	if currCountry == prevCountry {
		return nil, nil
	}

	prevTSInt, err := strconv.ParseInt(prevTSStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp in Redis: %v", err)
	}

	delta := nowTS - prevTSInt
	if delta > int64(d.window.Seconds()) {
		return nil, nil
	}

	alert := model.NewAlert(
		model.AnomalyGeoSwitching,
		[]*model.Event{e},
		model.AlertWarning,
		time.Now(),
		model.GeoSwitchingData{
			FromCountry: prevCountry,
			ToCountry:   currCountry,
			FromEventID: prevID,
			ToEventID:   e.ID,
			IntervalSec: delta,
		},
	)

	return alert, nil
}
