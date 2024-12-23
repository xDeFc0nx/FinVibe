package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func CreateWebSocketConnection(userID string) (string, error) {
	socket := new(types.WebSocketConnection)
	socket.ID = uuid.New().String()
	socket.UserID = userID
	socket.ConnectionID = uuid.New().String()
	socket.IsActive = true
	socket.LastPing = time.Now()
	socket.CreatedAt = time.Now()

	// Save to the database
	if err := db.DB.Create(socket).Error; err != nil {
		return "", fmt.Errorf("failed to create WebSocket: %w", err)
	}

	return socket.ID, nil
}
