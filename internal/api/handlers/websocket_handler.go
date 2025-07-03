package handlers

import (
	"net/http"

	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/websocket"
	"github.com/google/uuid"
	gorillaws "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketHandler struct {
	manager *websocket.Manager
}

func NewWebSocketHandler(manager *websocket.Manager) *WebSocketHandler {
	return &WebSocketHandler{
		manager: manager,
	}
}

var upgrader = gorillaws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for this example. In production, you should restrict this.
	},
}

// HandleWebSocket godoc
// @Summary Connect to WebSocket
// @Description Upgrade HTTP connection to WebSocket for real-time updates
// @Tags websocket
// @Accept json
// @Produce json
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /ws [get]
// @Security BearerAuth
func (h *WebSocketHandler) HandleWebSocket(c echo.Context) error {
	// Get user ID from JWT token
	userID := c.Get("user_id").(uint)

	// Upgrade connection
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	// Create new client
	client := &websocket.Client{
		ID:      uuid.New().String(),
		UserID:  userID,
		Conn:    conn,
		Manager: h.manager,
	}

	// Register client
	h.manager.Register <- client

	// Start listening for messages from this client
	go func() {
		defer func() {
			client.Conn.Close()
		}()

		for {
			_, _, err := client.Conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()

	return nil
}
