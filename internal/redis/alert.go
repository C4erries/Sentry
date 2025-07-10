package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/c4erries/Sentry/internal/model"
)

const (
	alertCacheTTL = 10 * time.Minute
)

type AlertCache interface {
	SaveAlert(ctx context.Context, alert *model.Alert) error
	GetAlert(ctx context.Context, alertID string) (*model.Alert, error)
}

func (r *Adapter) SaveAlert(ctx context.Context, alert *model.Alert) error {
	data, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("marshal alert: %w", err)
	}

	key := fmt.Sprintf("alerts:%s", alert.ID)
	return r.client.Set(ctx, key, data, alertCacheTTL).Err()
}

func (r *Adapter) GetAlert(ctx context.Context, alertID string) (*model.Alert, error) {
	key := fmt.Sprintf("alerts:%s", alertID)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("get alert: %w", err)
	}

	var alert model.Alert
	if err := json.Unmarshal([]byte(val), &alert); err != nil {
		return nil, fmt.Errorf("unmarshal alert: %w", err)
	}

	return &alert, nil
}
