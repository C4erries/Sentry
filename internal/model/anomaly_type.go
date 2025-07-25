package model

import (
	"encoding/json"
	"fmt"
)

type AnomalyType string

const (
	AnomalyLoginStorm   = "login_storm"
	AnomalyGeoSwitching = "geo_switching"
)

var validAnomalyTypes = []AnomalyType{
	AnomalyLoginStorm,
	AnomalyGeoSwitching,
}

func (a AnomalyType) IsValid() bool {
	for _, v := range validAnomalyTypes {
		if a == v {
			return true
		}
	}
	return false
}

func ParseAnomalyType(s string) (AnomalyType, error) {
	a := AnomalyType(s)
	if !a.IsValid() {
		return "", fmt.Errorf("invalid AnomalyType %q", s)
	}
	return a, nil
}

func (anomalyType AnomalyType) UnmarshalData(raw []byte) (interface{}, error) {
	switch anomalyType {
	case AnomalyLoginStorm:
		var d LoginStormData
		if err := json.Unmarshal(raw, &d); err != nil {
			return nil, err
		}
		return d, nil
	case AnomalyGeoSwitching:
		var d GeoSwitchingData
		if err := json.Unmarshal(raw, &d); err != nil {
			return nil, err
		}
		return d, nil
	default:
		return nil, fmt.Errorf("unknown anomaly type: %v", anomalyType)
	}
}

func (anomalyType AnomalyType) String() string {
	return string(anomalyType)
}
