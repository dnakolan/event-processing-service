package storage

import (
	"context"
	"testing"
	"time"

	"github.com/dnakolan/event-processing-service/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventStorage_Save(t *testing.T) {
	storage := NewEventStorage()
	ctx := context.Background()
	zeroTime := time.Unix(0, 0)
	uid := uuid.New().String()

	event := &models.Event{
		EventID:   uid,
		UserID:    "123",
		EventType: models.EventTypePageView,
		Timestamp: &zeroTime,
		Properties: models.EventProperties{
			Page:      "/home",
			Amount:    29.99,
			ProductID: "xyz",
		},
	}

	err := storage.Save(ctx, event)
	require.NoError(t, err)

	// Verify the event was saved
	saved, err := storage.FindById(ctx, event.EventID)
	require.NoError(t, err)
	assert.Equal(t, *event, *saved)
}

func TestEventStorage_FindById(t *testing.T) {
	storage := NewEventStorage()
	ctx := context.Background()
	notFoundUID := uuid.New()
	zeroTime := time.Unix(0, 0)
	uid := uuid.New().String()

	event := &models.Event{
		EventID:   uid,
		UserID:    "123",
		EventType: models.EventTypePageView,
		Timestamp: &zeroTime,
		Properties: models.EventProperties{
			Page:      "/home",
			Amount:    29.99,
			ProductID: "xyz",
		},
	}

	// Save a event first
	err := storage.Save(ctx, event)
	require.NoError(t, err)

	tests := []struct {
		name          string
		uid           string
		expectError   bool
		expectedError string
	}{
		{
			name:        "successful retrieval",
			uid:         event.EventID,
			expectError: false,
		},
		{
			name:          "event not found",
			uid:           notFoundUID.String(),
			expectError:   true,
			expectedError: "event not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := storage.FindById(ctx, tt.uid)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, found)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, found)
				assert.Equal(t, *event, *found)
			}
		})
	}
}

func TestEventStorage_FindAll(t *testing.T) {
	storage := NewEventStorage()
	ctx := context.Background()
	zeroTime := time.Unix(0, 0)
	midTime := time.Unix(50, 0)
	endTime := time.Unix(100, 0)

	// Create test events
	events := []*models.Event{
		&models.Event{
			UserID:    "123",
			EventType: models.EventTypePageView,
			Properties: models.EventProperties{
				Page:      "/home",
				Amount:    29.99,
				ProductID: "abc",
			},
		},
		&models.Event{
			EventID:   uuid.New().String(),
			UserID:    "123",
			EventType: models.EventTypePageView,
			Properties: models.EventProperties{
				Page:      "/about",
				Amount:    19.99,
				ProductID: "def",
			},
		},
		&models.Event{
			EventID:   uuid.New().String(),
			UserID:    "789",
			EventType: models.EventTypePageView,
			Properties: models.EventProperties{
				Page:      "/contact",
				Amount:    9.99,
				ProductID: "ghi",
			},
		},
	}

	// Set CreatedAt for all events
	for _, event := range events {
		event.Timestamp = &midTime
	}

	// Save all events
	for _, w := range events {
		err := storage.Save(ctx, w)
		require.NoError(t, err)
	}

	tests := []struct {
		name            string
		filter          *models.EventFilter
		expectedCount   int
		expectedUserIds []string
		expectedFilter  func(*models.Event) bool
	}{
		{
			name:            "no filter",
			filter:          nil,
			expectedCount:   3,
			expectedUserIds: []string{"123", "123", "789"},
		},
		{
			name: "filter by user id",
			filter: &models.EventFilter{
				UserID: stringPtr("123"),
			},
			expectedCount:   2,
			expectedUserIds: []string{"123", "123"},
			expectedFilter: func(w *models.Event) bool {
				return w.UserID == "123"
			},
		},
		{
			name: "filter by time",
			filter: &models.EventFilter{
				StartTimestamp: &zeroTime,
				EndTimestamp:   &endTime,
			},
			expectedCount:   3,
			expectedUserIds: []string{"123", "123", "789"},
			expectedFilter: func(w *models.Event) bool {
				return w.Timestamp.After(zeroTime) && w.Timestamp.Before(endTime)
			},
		},
		{
			name: "filter with no matches",
			filter: &models.EventFilter{
				UserID: stringPtr("nonexistent"),
			},
			expectedCount:   0,
			expectedUserIds: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := storage.FindAll(ctx, tt.filter)
			require.NoError(t, err)
			assert.Len(t, found, tt.expectedCount)

			// Verify the names of returned events
			names := make([]string, len(found))
			for i, w := range found {
				names[i] = w.UserID
			}
			assert.ElementsMatch(t, tt.expectedUserIds, names)

			// If there's a specific filter function, verify each event matches it
			if tt.expectedFilter != nil {
				for _, w := range found {
					assert.True(t, tt.expectedFilter(w))
				}
			}
		})
	}
}

// Helper functions to create pointers
func float64Ptr(v float64) *float64 {
	return &v
}

func stringPtr(v string) *string {
	return &v
}
