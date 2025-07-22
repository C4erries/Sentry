package storage

import (
	"fmt"
	"time"

	"github.com/c4erries/Sentry/internal/model"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AlertDB struct {
	ID         string         `gorm:"column:id;primaryKey"`
	Rule       string         `gorm:"column:rule"`
	Level      string         `gorm:"column:level"`
	DetectedAt time.Time      `gorm:"column:detected_at"`
	Data       datatypes.JSON `gorm:"column:data"`

	EventID string  `gorm:"column:event_id"`
	Event   EventDB `gorm:"foreignKey:EventID;references:ID"`
}

func (AlertDB) TableName() string {
	return "alerts"
}

type EventDB struct {
	ID         string         `gorm:"column:id;primaryKey"`
	UserID     string         `gorm:"column:user_id"`
	EventType  string         `gorm:"column:event_type"`
	CreatedAt  time.Time      `gorm:"column:created_at"`
	IP         string         `gorm:"column:ip"`
	GeoCountry string         `gorm:"column:geo_country"`
	Data       datatypes.JSON `gorm:"column:data"`
}

func (EventDB) TableName() string {
	return "events"
}

func (a *AlertDB) ToAlert() (model.Alert, error) {
	var err error

	event, err := a.Event.ToEvent()
	if err != nil {
		return model.Alert{}, fmt.Errorf("can't convert that eventDB to event: %v", err)
	}
	id, err := uuid.Parse(a.ID)
	if err != nil {
		return model.Alert{}, fmt.Errorf("uuid parsing error: %v", err)
	}
	alert := model.Alert{
		ID:         id,
		DetectedAt: a.DetectedAt,
		Event:      &event,
	}

	alert.Rule, err = model.ParseAnomalyType(a.Rule)
	if err != nil {
		return model.Alert{}, err
	}

	alert.Data, err = alert.Rule.UnmarshalData(a.Data)
	if err != nil {
		return model.Alert{}, err
	}

	alert.Level, err = model.ParseAlertLevel(a.Level)
	if err != nil {
		return model.Alert{}, err
	}

	return alert, nil
}

func (e *EventDB) ToEvent() (model.Event, error) {
	eventType, err := model.ParseEventType(e.EventType)
	if err != nil {
		return model.Event{}, fmt.Errorf("can't parse event type: %w", err)
	}

	data, err := eventType.UnmarshalData(e.Data)
	if err != nil {
		return model.Event{}, fmt.Errorf("can't unmarshal data: %w", err)
	}
	id, err := uuid.Parse(e.ID)
	if err != nil {
		return model.Event{}, fmt.Errorf("uuid parsing error: %w", err)
	}
	event := model.Event{
		ID: id,
		BaseEvent: model.BaseEvent{
			EventType:  eventType,
			UserID:     e.UserID,
			CreatedAt:  e.CreatedAt,
			IP:         e.IP,
			GeoCountry: e.GeoCountry,
		},
		Data: data,
	}

	return event, nil
}
