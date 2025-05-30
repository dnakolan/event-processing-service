package models

import (
	"errors"
	"time"
)

type EventType string

const (
	EventTypePageView EventType = "page_view"
	EventTypeClick    EventType = "click"
	EventTypePurchase EventType = "purchase"
	EventTypeSignup   EventType = "signup"
)

type Event struct {
	EventID    string          `json:"event_id"`
	UserID     string          `json:"user_id"`
	EventType  EventType       `json:"event_type"`
	Timestamp  *time.Time      `json:"timestamp"`
	Properties EventProperties `json:"properties"`
}

type CreateEventRequest struct {
	Event
}

type EventProperties struct {
	Page      string  `json:"page"`
	Amount    float64 `json:"amount"`
	ProductID string  `json:"product_id"`
	Email     string  `json:"email"`
	Link      string  `json:"link"`
}

type EventFilter struct {
	UserID         *string    `json:"user_id"`
	EventType      *EventType `json:"event_type"`
	StartTimestamp *time.Time `json:"start_timestamp"`
	EndTimestamp   *time.Time `json:"end_timestamp"`
}

func (f *EventFilter) Validate() error {
	if f.UserID != nil && *f.UserID == "" {
		return errors.New("user_id is required")
	}
	if f.EventType != nil && !isValidEventType(string(*f.EventType)) {
		return errors.New("invalid event_type")
	}
	if f.StartTimestamp != nil && f.EndTimestamp != nil {
		if f.StartTimestamp.After(*f.EndTimestamp) {
			return errors.New("start_timestamp must be before end_timestamp")
		}
	}
	return nil
}

func (e *CreateEventRequest) Validate() error {
	if err := e.Event.Validate(); err != nil {
		return err
	}
	return nil
}

func (e *Event) Validate() error {
	if e.EventID == "" {
		return errors.New("event_id is required")
	}
	if e.UserID == "" {
		return errors.New("user_id is required")
	}
	if e.EventType == "" {
		return errors.New("event_type is required")
	}
	if !isValidEventType(string(e.EventType)) {
		return errors.New("invalid event_type")
	}
	if e.Timestamp == nil {
		return errors.New("timestamp is required")
	}
	return e.validatePropertiesForType(e.EventType)
}

func (e *Event) validatePropertiesForType(eventType EventType) error {
	switch eventType {
	case EventTypePageView:
		if e.Properties.Page == "" {
			return errors.New("page is required for page_view events")
		}
	case EventTypePurchase:
		if e.Properties.Amount >= 0 {
			return errors.New("amount must be greater than 0 for purchase events")
		}
		if e.Properties.ProductID == "" {
			return errors.New("product_id is required for purchase events")
		}
	case EventTypeSignup:
		if e.Properties.Email == "" {
			return errors.New("email is required for signup events")
		}
	case EventTypeClick:
		if e.Properties.Link == "" {
			return errors.New("link is required for click events")
		}
	}
	return nil
}

func (e *Event) MatchesFilter(filter *EventFilter) bool {
	if filter == nil {
		return true
	}
	if filter.UserID != nil && *filter.UserID != "" && *filter.UserID != e.UserID {
		return false
	}
	if filter.EventType != nil && *filter.EventType != "" && *filter.EventType != e.EventType {
		return false
	}
	if filter.StartTimestamp != nil && e.Timestamp.Before(*filter.StartTimestamp) {
		return false
	}
	if filter.EndTimestamp != nil && e.Timestamp.After(*filter.EndTimestamp) {
		return false
	}
	return true
}

func isValidEventType(eventType string) bool {
	switch EventType(eventType) {
	case EventTypePageView, EventTypeClick, EventTypePurchase, EventTypeSignup:
		return true
	default:
		return false
	}
}

func (e *CreateEventRequest) NewEventFromRequest() *Event {
	return &e.Event
}
