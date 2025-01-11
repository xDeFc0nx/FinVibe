package handlers

import (
	"encoding/json"
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

func HeartBeat(ws *websocket.Conn, data json.RawMessage, userID string) {
	socket := new(types.WebSocketConnection)

	if err := db.DB.Where("user_id = ?", userID).Find(&socket).Error; err != nil {

		Message(ws, "Error: Could not find user")
		return
	}

	if err := db.DB.Model(&socket).Update("LastPing", time.Now().UTC()).Error; err != nil {
		Message(ws, "Error: Failed to update")
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
			return c.Status(500).
				JSON(fiber.Map{"error": "Failed to update WebSocket connection"})
		}
		return c.Status(101).
			JSON(fiber.Map{"message": "WebSocket connection updated successfully"})
	}
	return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
}

func HandleWebSocketConnection(ws *fiber.Ctx) error {
	log.Println("Incoming WebSocket connection request")

	socket := new(types.WebSocketConnection)
	token := ws.Cookies("jwt-token")

	if token == "" {
		log.Println("Token is missing")
		return ws.Status(400).JSON(fiber.Map{"error": "Token is missing"})
	}

	// Decode and validate the JWT token
	userID, connectionID, err := DecodeJWTToken(token)
	if err != nil {
		log.Printf("Failed to decode JWT token: %v\n", err)
		return ws.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	log.Printf("Decoded user ID from token: %s\n", userID)
	log.Printf("Decoded connectionID from token: %s\n", connectionID)

	if err := db.DB.Where("connection_id = ?", connectionID).First(socket).Error; err != nil {
		log.Printf("WebSocket not found for ID: %s\n", socket.ConnectionID)
		return ws.Status(404).JSON(fiber.Map{"error": "WebSocket not found"})
	}

	if !socket.IsActive {
		log.Printf("WebSocket with ID %s is not active\n", socket.ConnectionID)

		return ws.Status(403).
			JSON(fiber.Map{"error": "WebSocket is not active"})
	}

	return websocket.New(func(ws *websocket.Conn) {
		go func() {
			ticker := time.NewTicker(15 * time.Second)
			defer ticker.Stop()

			for {
				timeout := time.After(20 * time.Second)

				select {
				case <-ticker.C:
					if time.Since(socket.LastPing) > 15*time.Second {
						Message(ws, "ping")
						return
					}

				case <-timeout:
					if time.Since(socket.LastPing) > 20*time.Second {
						Message(ws, "Connection Timeout")

						ws.Close()
						return
					}
				}
			}
		}()

		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				logger.Debug("read: %s", err)
				break
			}

			logger.Debug("recv: %s", msg)

			var message struct {
				Action string          `json:"action"`
				Data   json.RawMessage `json:"data"`
			}

			if err := json.Unmarshal(msg, &message); err != nil {
				Message(ws, InvalidData)
				continue
			}

			if message.Action == "" {
				Message(ws, "Error: Action required")
				continue
			}

			handlersMap := map[string]func(ws *websocket.Conn, data json.RawMessage, userID string){
				"createAccount":      CreateAccount,
				"createBudget":       CreateBudget,
				"createGoal":         CreateGoal,
				"deleteAccount":      DeleteAccount,
				"deleteBudget":       DeleteBudget,
				"deleteGoal":         DeleteGoal,
				"deleteTransaction":  DeleteTransaction,
				"deleteUser":         DeleteUser,
				"getAccounts":        GetAccounts,
				"getBudgets":         GetBudgets,
				"getGoals":           GetGoals,
				"getTransactionById": GetTransactionById,
				"getTransactions":    GetTransactions,
				"getUser":            GetUser,
				"pong":               HeartBeat,
				"updateAccount":      UpdateAccount,
				"updateBudget":       UpdateBudget,
				"updateGoals":        UpdateGoal,
				"updateTransaction":  UpdateTransaction,
				"updateUser":         UpdateUser,
			}

			handler, exists := handlersMap[message.Action]
			if exists {
				handler(ws, message.Data, userID)
			} else {
				Message(ws, "Error: Unknown Action")
			}
		}

		defer func() {
			ws.Close()
		}()
	})(ws)
}
