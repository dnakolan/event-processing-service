package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/dnakolan/event-processing-service/internal/models"
)

type EventStorage interface {
	Save(ctx context.Context, Event *models.Event) error
	FindAll(ctx context.Context, filter *models.EventFilter) ([]*models.Event, error)
	FindById(ctx context.Context, uid string) (*models.Event, error)
	Delete(ctx context.Context, uid string) error
	Clear(ctx context.Context) error
}

type eventStorage struct {
	sync.RWMutex
	data map[string]*models.Event
}

func NewEventStorage() *eventStorage {
	return &eventStorage{
		data: make(map[string]*models.Event),
	}
}

func (s *eventStorage) Save(ctx context.Context, Event *models.Event) error {
	s.Lock()
	defer s.Unlock()
	s.data[Event.EventID] = Event
	return nil
}

func (s *eventStorage) FindAll(ctx context.Context, filter *models.EventFilter) ([]*models.Event, error) {
	s.RLock()
	defer s.RUnlock()
	Events := make([]*models.Event, 0, len(s.data))
	for _, Event := range s.data {
		if filter == nil || Event.MatchesFilter(filter) {
			Events = append(Events, Event)
		}
	}
	return Events, nil
}

func (s *eventStorage) FindById(ctx context.Context, uid string) (*models.Event, error) {
	s.RLock()
	defer s.RUnlock()
	Event, ok := s.data[uid]
	if !ok {
		return nil, errors.New("event not found")
	}
	return Event, nil
}

func (s *eventStorage) Delete(ctx context.Context, uid string) error {
	s.Lock()
	defer s.Unlock()
	delete(s.data, uid)
	return nil
}

func (s *eventStorage) Clear(ctx context.Context) error {
	s.Lock()
	defer s.Unlock()
	s.data = make(map[string]*models.Event)
	return nil
}
