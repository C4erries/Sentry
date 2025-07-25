package model

import (
	"time"

	"github.com/google/uuid"
)

type Alert struct {
	ID         uuid.UUID   `json:"id"`
	Rule       AnomalyType `json:"rule"`
	Event      *Event      `json:"event"`
	Level      AlertLevel  `json:"level"`
	DetectedAt time.Time   `json:"detected_at"`
	Data       interface{} `json:"data"`
}

func NewAlert(rule AnomalyType, event *Event, level AlertLevel, detectedAt time.Time, data interface{}) *Alert {
	return &Alert{
		Rule:       rule,
		Event:      event,
		Level:      level,
		DetectedAt: detectedAt,
		ID:         uuid.New(),
		Data:       data,
	}
}

type LoginStormData struct {
	EventIDs []string `json:"event_ids"`
}

type GeoSwitchingData struct {
	FromCountry string `json:"from_country"`
	ToCountry   string `json:"to_country"`
	FromEventID string `json:"from_event_id"`
	ToEventID   string `json:"to_event_id"`
	IntervalSec int64  `json:"interval_sec"`
}
