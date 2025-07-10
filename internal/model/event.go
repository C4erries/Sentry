package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/language"
)

type BaseEvent struct {
	EventType  EventType `json:"event_type"`
	UserId     string    `json:"user_id"`
	Timestamp  time.Time `json:"timestamp"`
	IP         string    `json:"ip"`
	GeoCountry string    `json:"geo_country"`
	Device     string    `json:"device"`
}

type Event struct {
	BaseEvent
	ID   string      `json:"id"` //uuid
	Data interface{} `json:"data"`
}

func NewEvent(baseEvent BaseEvent, data interface{}) (*Event, error) {
	event := &Event{
		BaseEvent: baseEvent,
		ID:        uuid.New().String(),
		Data:      data,
	}

	if err := event.Normalize(); err != nil {
		return nil, err
	}

	return event, nil
}

func (e *Event) Validate() error {
	if _, err := uuid.Parse(e.ID); err != nil {
		return fmt.Errorf("invalid UUID: %v", err)
	}
	if err := e.BaseEvent.EventType.Validate(); err != nil {
		return fmt.Errorf("invalid EventType: %v", err)
	}

	country, err := language.ParseRegion(strings.ToUpper(e.GeoCountry))
	if err != nil {
		return fmt.Errorf("CountryParse error: %v", err)
	}
	if !country.IsCountry() {
		e.GeoCountry = ""
		return fmt.Errorf("invalid country code: %v", country)
	}

	return nil
}

func (e *Event) Normalize() error {
	if e.GeoCountry != "" {
		country, err := language.ParseRegion(strings.ToUpper(e.GeoCountry))
		//log.Printf("Normalization [event-%s] Country: %s, err: %v", e.ID, country, err)
		if err != nil {
			e.GeoCountry = ""
			return fmt.Errorf("CountryParse error: %v", err)
		}
		if !country.IsCountry() {
			e.GeoCountry = ""
			return fmt.Errorf("invalid country code: %v", country)
		}
		e.GeoCountry = country.String()
	}

	return nil
}

/*
	func (e Event) GetGeo() (*string, error) {
		ip, err := netip.ParseAddr(e.IP)
		if err != nil {
			return nil, fmt.Errorf("invalid IP: %v", err)
		}
		db, err := geoip2.Open("GeoIP2-City.mmdb")
		if err != nil {
			return nil, fmt.Errorf("can't connect to geoip2: %v", err)
		}
		defer db.Close()

		record, err := db.City(ip)
		if err != nil {
			return nil, fmt.Errorf("db.City error:%v", err)
		}
		if !record.HasData() {
			return nil, fmt.Errorf("no data found for this IP")
		}
		return &record.Country.ISOCode, nil
	}
*/
type LoginData struct {
	Method  string `json:"method"`
	Success bool   `json:"success"`
}

type TransactionData struct {
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	PaymentMethod string  `json:"payment_method"`
}
