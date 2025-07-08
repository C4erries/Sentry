package model

import "fmt"

type EventType string

const (
	EventLogin       EventType = "login"
	EventTransaction EventType = "transaction"
)

func (e EventType) Validate() error {
	switch e {
	case EventLogin, EventTransaction:
		return nil
	}
	return fmt.Errorf("unknown event type")
}

func (e EventType) String() string {
	return string(e)
}
