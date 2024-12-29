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

func HeartBeat(c *websocket.Conn, data json.RawMessage, userID string) {
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

func HandleWebSocketConnection(c *fiber.Ctx) error {
	connectionID := c.Query("ConnectionID")
	socket := new(types.WebSocketConnection)

	if connectionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ConnectionID is required"})

	}

	if err := db.DB.Where("id = ?", connectionID).First(socket).Error; err != nil {

		return c.Status(404).JSON(fiber.Map{"error": "WebSocket not found"})

	}

	userID := socket.UserID

	if err := HandleCheckAuth(c, userID); err != nil {
		logger.Error("%s", err.Error())
	}

	if !socket.IsActive {
		return c.Status(403).JSON(fiber.Map{"error": "WebSocket is not active"})
	}

	return websocket.New(func(c *websocket.Conn) {
		go func() {
			ticker := time.NewTicker(15 * time.Second)
			defer ticker.Stop()

			for {
				timeout := time.After(5 * time.Second)

				select {
				case <-ticker.C:

					if time.Since(socket.LastPing) > 15*time.Second {
						if err := c.WriteMessage(websocket.TextMessage, []byte(`ping`)); err != nil {
							log.Printf("Error sending ping: %v", err)
							return
						}

					}

				case <-timeout:
					if time.Since(socket.LastPing) > 20*time.Second {
						log.Println("No response from client, closing connection")
						if err := c.WriteMessage(websocket.TextMessage, []byte(`connection timeout`)); err != nil {
							log.Printf("Error sending ping: %v", err)
							return
						}
						c.Close()
						return
					}
				}
			}

		}()
		var (
			msg []byte
			err error
		)
		for {
			if _, msg, err = c.ReadMessage(); err != nil {
				logger.Debug("read: %s", err)
				break
			}
			logger.Debug("recv: %s", msg)

			var message struct {
				Action string          `json:"action"`
				Data   json.RawMessage `json:"data"`
			}
			if err := json.Unmarshal(msg, &message); err != nil {
				if err := c.WriteMessage(websocket.TextMessage, []byte(`{"error":"Invalid message format"}`+err.Error())); err != nil {
					logger.Error("%s", err.Error())
				}
				continue
			}

			if message.Action == "" {
				if err := c.WriteMessage(websocket.TextMessage, []byte(`{"error":"Action is required"}`)); err != nil {
					logger.Error("%s", err.Error())
				}
				continue
			}
			handlersMap := map[string]func(c *websocket.Conn, data json.RawMessage, userID string){
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
			handler, exists := handlersMap[message.Action]
			if exists {
				handler(c, message.Data, userID)
			} else {
				if err := c.WriteMessage(websocket.TextMessage, []byte(`{"error":"Unknown action"}`)); err != nil {
					logger.Error("%s", err.Error())
				}
			}

		}
		defer func() {
			socket.IsActive = false
			db.DB.Save(socket)
			c.Close()

		}()

	})(c)

}
