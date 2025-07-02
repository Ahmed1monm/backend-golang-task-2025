package websocket

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

// EventType represents the type of WebSocket event
type EventType string

// Client represents a WebSocket client connection
type Client struct {
	ID       string
	UserID   uint
	Conn     *websocket.Conn
	Manager  *Manager
	mu       sync.Mutex
	isClosed bool
}

// Event represents a WebSocket event
type Event struct {
	Type    EventType   `json:"type"`
	Payload interface{} `json:"payload"`
}

// Manager handles WebSocket connections and broadcasting
type Manager struct {
	clients    map[*Client]bool
	userConns  map[uint][]*Client // Map user ID to their connections
	Register   chan *Client       // Channel for registering new clients
	unregister chan *Client
	broadcast  chan Event
	mu         sync.RWMutex
}

// NewManager creates a new WebSocket manager
func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*Client]bool),
		userConns:  make(map[uint][]*Client),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Event, 100), // Buffered channel to prevent blocking
	}
}

// Start starts the WebSocket manager
func (m *Manager) Start() {
	for {
		select {
		case client := <-m.Register:
			m.mu.Lock()
			m.clients[client] = true
			m.userConns[client.UserID] = append(m.userConns[client.UserID], client)
			m.mu.Unlock()

		case client := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				m.removeUserConnection(client)
				client.Conn.Close()
			}
			m.mu.Unlock()

		case event := <-m.broadcast:
			m.broadcastEvent(event)
		}
	}
}

// SendToUser sends an event to a specific user's connections
func (m *Manager) SendToUser(userID uint, event Event) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	clients := m.userConns[userID]
	for _, client := range clients {
		client.sendEvent(event)
	}
}

// Broadcast sends an event to all connected clients
func (m *Manager) Broadcast(event Event) {
	m.broadcast <- event
}

func (m *Manager) broadcastEvent(event Event) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for client := range m.clients {
		client.sendEvent(event)
	}
}

func (m *Manager) removeUserConnection(client *Client) {
	conns := m.userConns[client.UserID]
	for i, c := range conns {
		if c == client {
			m.userConns[client.UserID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	if len(m.userConns[client.UserID]) == 0 {
		delete(m.userConns, client.UserID)
	}
}

func (c *Client) sendEvent(event Event) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isClosed {
		return
	}

	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
		c.isClosed = true
		c.Manager.unregister <- c
	}
}
