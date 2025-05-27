package services

import (
	"context"

	"github.com/dnakolan/event-processing-service/internal/models"
	"github.com/dnakolan/event-processing-service/internal/storage"
)

type EventsService interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	GetEvent(ctx context.Context, id string) (*models.Event, error)
	GetEvents(ctx context.Context, filter *models.EventFilter) ([]*models.Event, error)
}

type eventsService struct {
	storage storage.EventStorage
}

func NewEventsService(storage storage.EventStorage) *eventsService {
	return &eventsService{storage: storage}
}

func (s *eventsService) CreateEvent(ctx context.Context, event *models.Event) error {
	return s.storage.Save(ctx, event)
}

func (s *eventsService) GetEvent(ctx context.Context, id string) (*models.Event, error) {
	return s.storage.FindById(ctx, id)
}

func (s *eventsService) GetEvents(ctx context.Context, filter *models.EventFilter) ([]*models.Event, error) {
	return s.storage.FindAll(ctx, filter)
}
