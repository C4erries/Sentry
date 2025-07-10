package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Alert struct {
	Rule       AnomalyType `json:"rule"`
	Events     []*Event    `json:"events"`
	Level      AlertLevel  `json:"level"`
	DetectedAt time.Time   `json:"detected_at"`
	ID         string      `json:"id"`
	Data       interface{} `json:"data"`
}

func NewAlert(rule AnomalyType, events []*Event, level AlertLevel, detectedAt time.Time, data interface{}) *Alert {
	return &Alert{
		Rule:       rule,
		Events:     events,
		Level:      level,
		DetectedAt: detectedAt,
		ID:         uuid.New().String(),
		Data:       data,
	}
}

type AnomalyType string

const (
	AnomalyLoginStorm   = "login_storm"
	AnomalyGeoSwitching = "geo_switching"
)

func (anomalyType AnomalyType) IsValid() bool {
	switch anomalyType {
	case AnomalyGeoSwitching, AnomalyLoginStorm:
		return true
	default:
		return false
	}
}

type AlertLevel string

const (
	AlertWarning  = "warning"
	AlertCritical = "critical"
)

func (alertLevel AlertLevel) IsValid() bool {
	switch alertLevel {
	case AlertWarning, AlertCritical:
		return true
	default:
		return false
	}
}

type LoginStormData struct {
	EventIDs []string `json:"event_ids"`
}

func (d LoginStormData) String() string {
	str :=
		`LoginStormData
			EventIDs:
		`
	for _, id := range d.EventIDs {
		str += fmt.Sprintf(" - %s\n", id)
	}
	str = "" //!!!
	return str
}

func (d LoginStormData) IsPrintable() bool {
	return false
}

type GeoSwitchingData struct {
	FromCountry string `json:"from_country"`
	ToCountry   string `json:"to_country"`
	FromEventID string `json:"from_event_id"`
	ToEventID   string `json:"to_event_id"`
	IntervalSec int64  `json:"interval_sec"`
}

func (d GeoSwitchingData) String() string {
	str := fmt.Sprintf(
		`GeoSwitchingData
FromCountry: %s,
ToCountry: %s,
FromEventID: %s,
ToEventID: %s,
IntervalSec: %v`,
		d.FromCountry, d.ToCountry,
		d.FromEventID, d.ToEventID,
		d.IntervalSec,
	)
	return str
}

func (d GeoSwitchingData) IsPrintable() bool {
	return true
}
