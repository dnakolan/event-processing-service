package connections

import (
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/dnakolan/event-processing-service/internal/models"
	"github.com/gorilla/websocket"
)

type ConnectionManager interface {
	AddConnection(conn *websocket.Conn)
	RemoveConnection(conn *websocket.Conn)
	BroadcastEvent(event *models.Event)
}

type connectionManager struct {
	connections map[*websocket.Conn]bool
	mutex       sync.RWMutex
}

func NewConnectionManager() *connectionManager {
	return &connectionManager{
		connections: make(map[*websocket.Conn]bool),
	}
}

func (cm *connectionManager) AddConnection(conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.connections[conn] = true
}

func (cm *connectionManager) RemoveConnection(conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	delete(cm.connections, conn)
}

func (cm *connectionManager) BroadcastEvent(event *models.Event) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	eventJSON, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal event", "error", err.Error())
		return
	}

	for conn := range cm.connections {
		if err := conn.WriteMessage(websocket.TextMessage, eventJSON); err != nil {
			slog.Error("broadcast failed", "error", err.Error())
			// Connection is broken, remove it
			go func(c *websocket.Conn) {
				cm.RemoveConnection(c)
				c.Close()
			}(conn)
		}
	}
}
