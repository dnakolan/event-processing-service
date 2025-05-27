package models

import "time"

type Analytics struct {
	TimeWindow    string            `json:"time_window"`
	TotalEvents   int               `json:"total_events"`
	EventsByType  map[EventType]int `json:"events_by_type"`
	UniqueUsers   int               `json:"unique_users"`
	EventsPerHour []EventPerHour    `json:"events_per_hour"`
}

type EventPerHour struct {
	Hour  time.Time `json:"hour"`
	Count int       `json:"count"`
}

/*
{
  "time_window": "24h",
  "total_events": 1250,
  "events_by_type": {
    "page_view": 800,
    "click": 350,
    "purchase": 100
  },
  "unique_users": 45,
  "events_per_hour": [
    {"hour": "2025-05-26T14:00:00Z", "count": 120},
    {"hour": "2025-05-26T15:00:00Z", "count": 95}
  ]
}
*/
