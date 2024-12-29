package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xDeFc0nx/logger-go-pkg"

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

func HeartBeat(c *websocket.Conn, userID string) {
	socket := new(types.WebSocketConnection)

	if err := db.DB.Where("user_id = ?", userID).Find(&socket).Error; err != nil {

		if err := c.WriteMessage(websocket.TextMessage, []byte(err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}
	socket.LastPing = time.Now().UTC()

	if err := db.DB.Save(socket).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
	}

}
func HandleCheckAuth(c *fiber.Ctx, userID string) error {
	if CheckAuth(c) == nil {

		socket := new(types.WebSocketConnection)

		if err := db.DB.Where("user_id = ?", userID).Find(&socket).Error; err != nil {

			fmt.Printf("WebSocket not found: %v", err)
			return c.Status(404).JSON(fiber.Map{"error": "WebSocket not found"})

		}
		socket.LastPing = time.Now()
		socket.IsActive = true
		if err := db.DB.Save(socket).Error; err != nil {
			log.Printf("Failed to update WebSocket connection: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update WebSocket connection"})
		}
		return c.Status(101).JSON(fiber.Map{"message": "WebSocket connection updated successfully"})
	}
	return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
}
