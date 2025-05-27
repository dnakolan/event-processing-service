package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dnakolan/event-processing-service/internal/connections"
	"github.com/dnakolan/event-processing-service/internal/models"
	"github.com/dnakolan/event-processing-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type EventsHandler struct {
	connections connections.ConnectionManager
	service     services.EventsService
	upgrader    *websocket.Upgrader
}

func NewEventsHandler(service services.EventsService) *EventsHandler {
	return &EventsHandler{
		service:     service,
		upgrader:    &websocket.Upgrader{},
		connections: connections.NewConnectionManager(),
	}
}

func (h *EventsHandler) CreateEventsHandler(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()
	h.connections.AddConnection(conn)
	defer h.connections.RemoveConnection(conn)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			slog.Error("failed to read message", "error", err.Error())
			break
		}
		switch messageType {
		case websocket.TextMessage:
			// Process your JSON events
			h.handleEventJSON(c.Request.Context(), message)
		case websocket.BinaryMessage:
			// Maybe reject or handle differently
			slog.Error("binary data not supported")
		case websocket.CloseMessage:
			slog.Info("client disconnecting")
			return
		}
	}
}

func (h *EventsHandler) handleEventJSON(ctx context.Context, message []byte) {
	var event models.Event
	if err := json.Unmarshal(message, &event); err != nil {
		slog.Error("failed to unmarshal event", "error", err.Error())
		return
	}
	h.service.CreateEvent(ctx, &event)
	h.connections.BroadcastEvent(&event)
}
