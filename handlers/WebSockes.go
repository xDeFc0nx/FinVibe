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

		if err := ws.WriteMessage(websocket.TextMessage, []byte(err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}
	socket.LastPing = time.Now().UTC()

	if err := db.DB.Save(socket).Error; err != nil {
		if err := ws.WriteMessage(websocket.TextMessage, []byte(err.Error())); err != nil {
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

func HandleWebSocketConnection(ws *fiber.Ctx) error {

	socket := new(types.WebSocketConnection)
	token := ws.Query("token")

	if token == "" {
		return ws.Status(400).JSON(fiber.Map{"error": "Token is missing"})
	}

	// Decode and validate the JWT token
	userID, err := DecodeJWTToken(token)
	if err != nil {
		return ws.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Fetch WebSocket connection by ID
	if socket.ConnectionID == "" {
		return ws.Status(400).JSON(fiber.Map{"error": "ConnectionID is required"})
	}

	if err := db.DB.Where("id = ?", socket.ConnectionID).First(socket).Error; err != nil {
		return ws.Status(404).JSON(fiber.Map{"error": "WebSocket not found"})
	}

	if !socket.IsActive {
		return ws.Status(403).JSON(fiber.Map{"error": "WebSocket is not active"})
	}

	// Accept WebSocket connection
	return websocket.New(func(ws *websocket.Conn) {
		// Connection ping/pong handling
		go func() {
			ticker := time.NewTicker(15 * time.Second)
			defer ticker.Stop()

			for {
				timeout := time.After(20 * time.Second)

				select {
				case <-ticker.C:
					if time.Since(socket.LastPing) > 15*time.Second {
						if err := ws.WriteMessage(websocket.TextMessage, []byte(`ping`)); err != nil {
							log.Printf("Error sending ping: %v", err)
							return
						}
					}

				case <-timeout:
					if time.Since(socket.LastPing) > 20*time.Second {
						log.Println("No response from client, closing connection")
						if err := ws.WriteMessage(websocket.TextMessage, []byte(`connection timeout`)); err != nil {
							log.Printf("Error sending ping: %v", err)
							return
						}
						ws.Close()
						return
					}
				}
			}
		}()

		// Handle WebSocket messages
		for {
			msgType, msg, err := ws.ReadMessage()
			if err != nil {
				logger.Debug("read: %s", err)
				break
			}

			logger.Debug("recv: %s", msg)

			var message struct {
				Action string          `json:"action"`
				Data   json.RawMessage `json:"data"`
			}

			// Handle invalid message format
			if err := json.Unmarshal(msg, &message); err != nil {
				if err := ws.WriteMessage(msgType, []byte(`{"error":"Invalid message format"}`+err.Error())); err != nil {
					logger.Error("%s", err.Error())
				}
				continue
			}

			// Handle missing action
			if message.Action == "" {
				if err := ws.WriteMessage(msgType, []byte(`{"error":"Action is required"}`)); err != nil {
					logger.Error("%s", err.Error())
				}
				continue
			}

			// Map actions to handler functions
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
				"logout":             LogoutHandler,
				"pong":               HeartBeat,
				"updateAccount":      UpdateAccount,
				"updateBudget":       UpdateBudget,
				"updateGoals":        UpdateGoal,
				"updateTransaction":  UpdateTransaction,
				"updateUser":         UpdateUser,
			}

			// Call the appropriate handler based on the action
			handler, exists := handlersMap[message.Action]
			if exists {
				handler(ws, message.Data, userID)
			} else {
				if err := ws.WriteMessage(msgType, []byte(`{"error":"Unknown action"}`)); err != nil {
					logger.Error("%s", err.Error())
				}
			}
		}

		// Deactivate WebSocket when done
		defer func() {
			socket.IsActive = false
			db.DB.Save(socket)
			ws.Close()
		}()
	})(ws)
}
