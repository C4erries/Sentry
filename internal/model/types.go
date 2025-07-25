package model

import (
	"encoding/json"
	"fmt"
)

type EventType string

const (
	EventLogin       EventType = "login"
	EventTransaction EventType = "transaction"
)

var validEventTypes = []EventType{
	EventLogin,
	EventTransaction,
}

func (eventType EventType) Validate() error {
	for _, valid := range validEventTypes {
		if eventType == valid {
			return nil
		}
	}
	return fmt.Errorf("unknown event type")
}

func (eventType EventType) String() string {
	return string(eventType)
}

func (eventType EventType) UnmarshalData(raw []byte) (interface{}, error) {
	switch eventType {
	case EventLogin:
		var d LoginData
		if err := json.Unmarshal(raw, &d); err != nil {
			return nil, err
		}
		return d, nil
	case EventTransaction:
		var d TransactionData
		if err := json.Unmarshal(raw, &d); err != nil {
			return nil, err
		}
		return d, nil
	default:
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}
}

func ParseEventType(str string) (EventType, error) {
	eventType := EventType(str)
	if err := eventType.Validate(); err != nil {
		return "", err
	}
	return eventType, nil
}
