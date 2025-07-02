package events

type EventType string

const (
	EventLogin       EventType = "login"
	EventTransaction EventType = "transaction"
)

func (e EventType) IsValid() bool {
	switch e {
	case EventLogin, EventTransaction:
		return true
	}
	return false
}
