package services

import (
	"context"
	"sort"
	"time"

	"github.com/dnakolan/event-processing-service/internal/models"
	"github.com/dnakolan/event-processing-service/internal/storage"
)

type AnalyticsService interface {
	GetAnalytics(ctx context.Context, filter *models.EventFilter) (*models.Analytics, error)
}

type analyticsService struct {
	storage storage.EventStorage
}

func NewAnalyticsService(storage storage.EventStorage) AnalyticsService {
	return &analyticsService{storage: storage}
}

func (s *analyticsService) GetAnalytics(ctx context.Context, filter *models.EventFilter) (*models.Analytics, error) {
	events, err := s.storage.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &models.Analytics{
		TimeWindow:    filter.EndTimestamp.Sub(*filter.StartTimestamp),
		TotalEvents:   len(events),
		EventsByType:  eventsByType(events),
		UniqueUsers:   uniqueUsers(events),
		EventsPerHour: eventsPerHour(events),
	}, nil
}

func eventsByType(events []*models.Event) map[models.EventType]int {
	eventsByType := make(map[models.EventType]int)
	for _, event := range events {
		eventsByType[event.EventType]++
	}
	return eventsByType
}

func uniqueUsers(events []*models.Event) int {
	uniqueUsers := make(map[string]bool)
	for _, event := range events {
		uniqueUsers[event.UserID] = true
	}
	return len(uniqueUsers)
}

func eventsPerHour(events []*models.Event) []models.EventPerHour {
	eventsPerHourMap := make(map[time.Time]int)
	for _, event := range events {
		hour := event.Timestamp.Truncate(time.Hour)
		eventsPerHourMap[hour]++
	}

	eventsPerHour := make([]models.EventPerHour, 0)
	for hour, count := range eventsPerHourMap {
		eventsPerHour = append(eventsPerHour, models.EventPerHour{
			Hour:  hour,
			Count: count,
		})
	}
	sort.Slice(eventsPerHour, func(i, j int) bool {
		return eventsPerHour[i].Hour.Before(eventsPerHour[j].Hour)
	})
	return eventsPerHour
}
