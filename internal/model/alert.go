package model

import (
	"time"
)

type Alert struct {
	Rule       AnomalyType
	Events     []Event
	Level      AlertLevel
	DetectedAt time.Time
	Data       interface{}
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
	EventIDs []string
}
