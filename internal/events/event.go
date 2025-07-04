package events

import "time"

type BaseEvent struct {
	EventType EventType `json:"event_type"`
	UserId    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	IP        string    `json:"ip"`
	Device    string    `json:"device"`
}

type Event struct {
	BaseEvent
	Data interface{} `json:"data"`
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
