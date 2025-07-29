package testutils

import (
	"encoding/json"
	"net/netip"
	"testing"
	"time"

	"github.com/c4erries/Sentry/internal/model"
	"github.com/google/uuid"
)

func ValidEvent(t *testing.T) *model.Event {
	t.Helper()

	validEvent, err := model.NewEvent(
		model.BaseEvent{
			EventType:  model.EventLogin,
			UserID:     uuid.New().String(),
			CreatedAt:  time.Now(),
			IP:         netip.IPv4Unspecified().String(),
			GeoCountry: "RU",
		},
		model.LoginData{Method: "oauth", Success: false},
	)
	if err != nil {
		t.Fatalf("can't create new valid event: %v", err)
		return nil
	}

	if err := validEvent.Validate(); err != nil {
		t.Fatalf("valid event validation failed: %v", err)
		return nil
	}

	return validEvent
}

func ValidEventRaw(t *testing.T) []byte {
	t.Helper()

	validEvent := ValidEvent(t)

	validEventData, err := json.Marshal(validEvent)
	if err != nil {
		t.Fatalf("failed to marshal valid event: %v", err)
		return nil
	}

	return validEventData
}

func InvalidEvent(t *testing.T) *model.Event {
	t.Helper()

	invalidEvent := &model.Event{
		BaseEvent: model.BaseEvent{
			EventType:  "unknown event type",
			UserID:     uuid.New().String(),
			CreatedAt:  time.Now(),
			IP:         netip.IPv4Unspecified().String(),
			GeoCountry: "RU",
		},
		Data: map[string]interface{}{"unkField": "31fds", "tpyea": map[string][]int{"213": {12, 23}}},
	}

	if invalidEvent.Validate() == nil {
		t.Fatalf("invalid event is valid")
		return nil
	}

	return invalidEvent
}

func InvalidEventRaw(t *testing.T) []byte {
	invalidEvent := InvalidEvent(t)

	invalidEventRaw, err := json.Marshal(invalidEvent)
	if err != nil {
		t.Fatalf("failed to marshal invalid event: %v", err)
		return nil
	}

	return invalidEventRaw
}
