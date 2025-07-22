package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/c4erries/Sentry/internal/model"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Save(ctx context.Context, event *model.Event) error {
	data, err := json.Marshal(event.Data)
	if err != nil {
		return fmt.Errorf("marshal event data: %w", err)
	}

	eventDB := EventDB{
		ID:         event.ID.String(),
		UserID:     event.UserID,
		EventType:  event.EventType.String(),
		CreatedAt:  event.CreatedAt,
		GeoCountry: event.GeoCountry,
		IP:         event.IP,
		Data:       datatypes.JSON(data),
	}

	if err := r.db.WithContext(ctx).Create(&eventDB).Error; err != nil {
		return fmt.Errorf("insert event: %w", err)
	}

	return nil
}
