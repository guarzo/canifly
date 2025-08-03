package handlers

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin (adjust for production)
		return true
	},
}

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// WebSocketHub manages WebSocket connections
type WebSocketHub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan WSMessage
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	logger     interfaces.Logger
	mu         sync.RWMutex
}

// NewWebSocketHub creates a new WebSocket hub
func NewWebSocketHub(logger interfaces.Logger) *WebSocketHub {
	return &WebSocketHub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan WSMessage),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		logger:     logger,
	}
}

// Run starts the WebSocket hub
func (h *WebSocketHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			h.logger.Infof("WebSocket client connected. Total: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()
			h.logger.Infof("WebSocket client disconnected. Total: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				if err := client.WriteJSON(message); err != nil {
					h.logger.Errorf("Error writing to WebSocket: %v", err)
					delete(h.clients, client)
					client.Close()
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastUpdate sends an update to all connected clients
func (h *WebSocketHub) BroadcastUpdate(updateType string, data interface{}) {
	select {
	case h.broadcast <- WSMessage{Type: updateType, Data: data}:
	default:
		h.logger.Warnf("WebSocket broadcast channel full, dropping message")
	}
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Errorf("WebSocket upgrade failed: %v", err)
		return
	}

	h.register <- conn

	// Keep connection alive and handle incoming messages
	go func() {
		defer func() {
			h.unregister <- conn
		}()

		for {
			var msg WSMessage
			if err := conn.ReadJSON(&msg); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					h.logger.Errorf("WebSocket error: %v", err)
				}
				break
			}

			// Handle incoming messages (ping, etc.)
			switch msg.Type {
			case "ping":
				if err := conn.WriteJSON(WSMessage{Type: "pong", Data: "pong"}); err != nil {
					h.logger.Errorf("Error sending pong: %v", err)
					break
				}
			}
		}
	}()
}
