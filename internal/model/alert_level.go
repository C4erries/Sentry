package model

import "fmt"

type AlertLevel string

const (
	AlertWarning  = "warning"
	AlertCritical = "critical"
)

var validAlertLevels = []AlertLevel{
	AlertWarning,
	AlertCritical,
}

func (alertLevel AlertLevel) IsValid() bool {
	for _, level := range validAlertLevels {
		if level == alertLevel {
			return true
		}
	}
	return false
}
func ParseAlertLevel(s string) (AlertLevel, error) {
	a := AlertLevel(s)
	if !a.IsValid() {
		return "", fmt.Errorf("invalid AlertType %q", s)
	}
	return a, nil
}
