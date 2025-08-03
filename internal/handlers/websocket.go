package handlers

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

// isAllowedOrigin checks if the origin is allowed
func isAllowedOrigin(origin string) bool {
	// Allow localhost connections on common development ports
	allowedOrigins := []string{
		"http://localhost:3000",  // React dev server
		"http://localhost:5173",  // Vite dev server
		"http://localhost:42423", // Backend server
		"http://127.0.0.1:3000",
		"http://127.0.0.1:5173",
		"http://127.0.0.1:42423",
		"app://.", // Electron app
		"file://", // Local file access (Electron)
	}

	// Check exact match
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}

	// For production, you would load allowed origins from configuration
	// For now, also allow empty origin (same-origin requests)
	return origin == ""
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return isAllowedOrigin(origin)
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
	done       chan struct{}
	running    bool
}

// NewWebSocketHub creates a new WebSocket hub
func NewWebSocketHub(logger interfaces.Logger) *WebSocketHub {
	return &WebSocketHub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan WSMessage),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		logger:     logger,
		done:       make(chan struct{}),
		running:    false,
	}
}

// Run starts the WebSocket hub
func (h *WebSocketHub) Run() {
	h.mu.Lock()
	h.running = true
	h.mu.Unlock()

	for {
		select {
		case <-h.done:
			h.logger.Info("WebSocket hub shutting down...")
			h.closeAllClients()
			return

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
			clientsCopy := make(map[*websocket.Conn]bool)
			for client := range h.clients {
				clientsCopy[client] = true
			}
			h.mu.RUnlock()

			// Write to clients without holding the read lock
			for client := range clientsCopy {
				if err := client.WriteJSON(message); err != nil {
					h.logger.Errorf("Error writing to WebSocket: %v", err)
					// Send to unregister channel instead of directly manipulating
					select {
					case h.unregister <- client:
					default:
						// Channel might be full, but that's okay as the client
						// will be cleaned up eventually
					}
				}
			}
		}
	}
}

// Shutdown gracefully shuts down the WebSocket hub
func (h *WebSocketHub) Shutdown() {
	h.mu.Lock()
	if h.running {
		close(h.done)
		h.running = false
	}
	h.mu.Unlock()
}

// closeAllClients closes all connected WebSocket clients
func (h *WebSocketHub) closeAllClients() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		client.Close()
	}
	h.clients = make(map[*websocket.Conn]bool)
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
