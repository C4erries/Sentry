package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/c4erries/Sentry/internal/model"
)

type AlertDB struct {
	ID         string          `db:"alert_id"`
	Rule       string          `db:"alert_rule"`
	Level      string          `db:"alert_level"`
	DetectedAt time.Time       `db:"alert_detected_at"`
	RawData    json.RawMessage `db:"alert_data"`
	Event      EventDB         `db:"-"`
}

type EventDB struct {
	ID         string          `db:"event_id"`
	UserID     string          `db:"event_user_id"`
	Timestamp  time.Time       `db:"event_timestamp"`
	EventType  string          `db:"event_type"`
	IP         string          `db:"event_ip"`
	GeoCountry string          `db:"event_geo_country"`
	RawData    json.RawMessage `db:"event_data"`
}

func (a *AlertDB) ToAlert() (model.Alert, error) {
	var err error

	event, err := a.Event.ToEvent()
	if err != nil {
		return model.Alert{}, fmt.Errorf("can't convert that eventDB to event: %v", err)
	}

	alert := model.Alert{
		ID:         a.ID,
		DetectedAt: a.DetectedAt,
		Event:      &event,
	}

	alert.Rule, err = model.ParseAnomalyType(a.Rule)
	if err != nil {
		return model.Alert{}, err
	}

	alert.Data, err = alert.Rule.UnmarshalData(a.RawData)
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

	data, err := eventType.UnmarshalData(e.RawData)
	if err != nil {
		return model.Event{}, fmt.Errorf("can't unmarshal data: %w", err)
	}

	event := model.Event{
		ID: e.ID,
		BaseEvent: model.BaseEvent{
			EventType:  eventType,
			UserID:     e.UserID,
			Timestamp:  e.Timestamp,
			IP:         e.IP,
			GeoCountry: e.GeoCountry,
		},
		Data: data,
	}

	return event, nil
}
