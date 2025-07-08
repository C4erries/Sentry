package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type BaseEvent struct {
	EventType EventType `json:"event_type"`
	UserId    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	IP        string    `json:"ip"`
	Device    string    `json:"device"`
}

type Event struct {
	BaseEvent
	ID   string      `json:"id"` //uuid
	Data interface{} `json:"data"`
}

func NewEvent(baseEvent BaseEvent, data interface{}) Event {
	return Event{
		BaseEvent: baseEvent,
		ID:        uuid.New().String(),
		Data:      data,
	}
}

func (e Event) Validate() error {
	if _, err := uuid.Parse(e.ID); err != nil {
		return fmt.Errorf("invalid UUID: %v", err)
	}
	if err := e.BaseEvent.EventType.Validate(); err != nil {
		return fmt.Errorf("invalid EventType: %v", err)
	}
	return nil
}

type LoginData struct {
	Method  string `json:"method"`
	Success bool   `json:"success"`
}

type TransactionData struct {
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	PaymentMethod string  `json:"payment_method"`
}
