package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/c4erries/Sentry/internal/model"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AlertRepository struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) *AlertRepository {
	return &AlertRepository{db: db}
}

func (r *AlertRepository) Save(ctx context.Context, alert *model.Alert) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		alertDB := AlertDB{
			ID:         alert.ID.String(),
			Rule:       alert.Rule.String(),
			Level:      alert.Level.String(),
			DetectedAt: alert.DetectedAt,
			EventID:    alert.Event.ID.String(),
		}

		// Сериализуем Data в JSON
		data, err := json.Marshal(alert.Data)
		if err != nil {
			return fmt.Errorf("marshal alert data: %w", err)
		}
		alertDB.Data = datatypes.JSON(data)

		// Сохраняем alert
		if err := tx.Create(&alertDB).Error; err != nil {
			return fmt.Errorf("insert alert: %w", err)
		}

		// Если нужно вставлять в alert_events, можно раскомментировать и адаптировать:
		/*
			alertEvent := AlertEvent{
				AlertID: alert.ID,
				EventID: alert.Event.ID,
			}
			if err := tx.Create(&alertEvent).Error; err != nil {
				return fmt.Errorf("insert alert_event: %w", err)
			}
		*/

		return nil
	})
}

func (r *AlertRepository) AllEventsExist(ctx context.Context, eventIDs []uuid.UUID) (bool, error) {
	if len(eventIDs) == 0 {
		return true, nil
	}

	var count int64
	err := r.db.WithContext(ctx).
		Model(&EventDB{}).
		Where("id IN ?", eventIDs).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("checking events exist: %w", err)
	}

	return count == int64(len(eventIDs)), nil
}

func (r *AlertRepository) FindByID(ctx context.Context, alertID string) (*AlertDB, error) {
	var alert AlertDB

	err := r.db.WithContext(ctx).
		Preload("Event").
		First(&alert, "id = ?", alertID).Error

	if err != nil {
		return nil, err
	}

	return &alert, nil
}

// TODO: json_agg
func (r *AlertRepository) FindByUser(ctx context.Context, userID string) ([]*model.Alert, error) {
	var alertDBs []AlertDB
	if err := r.db.
		WithContext(ctx).
		Preload("Event").
		Joins("JOIN events ON events.id = alerts.event_id").
		Where("events.user_id = ?", userID).
		Order("alerts.detected_at DESC").
		Find(&alertDBs).Error; err != nil {
		return nil, err
	}

	var result []*model.Alert
	var parseErrors *multierror.Error

	for _, raw := range alertDBs {
		alert, err := raw.ToAlert()
		if err != nil {
			parseErrors = multierror.Append(parseErrors, fmt.Errorf("failed to convert raw alert: %w", err))
			continue
		}
		result = append(result, &alert)
	}

	slog.Debug("[FindByUser]", slog.Int("count", len(result)), slog.String("userID", userID))
	return result, parseErrors.ErrorOrNil()
}
