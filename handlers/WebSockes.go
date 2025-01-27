package handlers

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
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

	if err := db.DB.Create(socket).Error; err != nil {
		slog.Error(
			"Creating WebSocket failed",
			slog.String("UserID", userID),
			slog.String("error", err.Error()),
		)
		return "", err
	}

	return socket.ID, nil
}

func HeartBeat(ws *websocket.Conn, data json.RawMessage, userID string) {
	socket := new(types.WebSocketConnection)

	if err := db.DB.Where("user_id = ?", userID).Find(&socket).Error; err != nil {

		Send_Error(ws, "Could not find user", err)
		return
	}

	if err := db.DB.Model(&socket).Update("LastPing", time.Now().UTC()).Error; err != nil {
		Send_Error(ws, "Failed to update", err)
	}
}

func HandleCheckAuth(c *fiber.Ctx, userID string) error {
	if CheckAuth(c) == nil {

		socket := new(types.WebSocketConnection)

		if err := db.DB.Where("user_id = ?", userID).Find(&socket).Error; err != nil {
			slog.Error(
				"WebSocket not found",
				slog.String("UserID", userID),
				slog.String("error", err.Error()),
			)

			return c.Status(404).JSON(fiber.Map{"error": "WebSocket not found"})

		}
		socket.LastPing = time.Now()
		socket.IsActive = true
		if err := db.DB.Save(socket).Error; err != nil {
			slog.Error(
				"WebSocket update failed",
				slog.String("UserID", userID),
				slog.String("error", err.Error()),
			)

			return c.Status(500).
				JSON(fiber.Map{"error": "Failed to update WebSocket connection"})
		}
		return c.Status(101).
			JSON(fiber.Map{"message": "WebSocket connection updated successfully"})
	}
	return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
}

func HandleWebSocketConnection(c *fiber.Ctx) error {
	slog.Info("Incoming Socket connection", slog.String("IP Address", c.IP()))

	socket := new(types.WebSocketConnection)
	token := c.Cookies("jwt-token")

	if token == "" {
		slog.Info("Token is missing", slog.String("error", "Token is missing"))
		return c.Status(400).JSON(fiber.Map{"error": "Token is missing"})
	}

	userID, connectionID, err := DecodeJWTToken(token)
	if err != nil {
		slog.Error(
			"Failed to decode JWT token",
			slog.String("error", err.Error()),
		)

		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := db.DB.Where("connection_id = ?", connectionID).First(socket).Error; err != nil {
		slog.Error(
			"WebSocket not found",
			slog.String("UserID", userID),
			slog.String("error", err.Error()),
		)

		return c.Status(404).JSON(fiber.Map{"error": "WebSocket not found"})
	}

	socket.IsActive = true
	if err := db.DB.Save(socket).Error; err != nil {
		slog.Error("Failed to update socket", slog.String("error", err.Error()))
	}
	return websocket.New(func(ws *websocket.Conn) {
		slog.Info(
			"User Connected",
			slog.String("UserID", userID),
		)
		go func() {
			ticker := time.NewTicker(15 * time.Second)
			defer ticker.Stop()

			for {
				timeout := time.After(20 * time.Second)

				select {
				case <-ticker.C:
					if socket.IsActive {
						if time.Since(socket.LastPing) > 15*time.Second {
							response := map[string]string{
								"message": "pong",
							}
							responseJson, _ := json.Marshal(response)
							Send_Message(ws, string(responseJson))
						}
					}

				case <-timeout:
					if socket.IsActive {
						if time.Since(socket.LastPing) > 20*time.Second {
							response := map[string]string{
								"message": "Timeout",
							}
							responseJson, _ := json.Marshal(response)
							Send_Message(ws, string(responseJson))
							ws.Close()
							return
						}
					}
				}
			}
		}()

		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				slog.Info("read", slog.String("error", err.Error()))
				break
			}
			var message struct {
				Action string          `json:"action"`
				Data   json.RawMessage `json:"data"`
			}

			if err := json.Unmarshal(msg, &message); err != nil {
				Send_Error(ws, InvalidData, err)
				continue
			}
			slog.Info(
				"recv",
				slog.String("action", message.Action),
				slog.String("from", userID),
			)

			if message.Action == "" {
				Send_Error(ws, "Action required", err)
				continue
			}

			handlersMap := map[string]func(ws *websocket.Conn, data json.RawMessage, userID string){
				"createAccount":      CreateAccount,
				"createBudget":       CreateBudget,
				"createGoal":         CreateGoal,
				"createTransaction":  CreateTransaction,
				"deleteAccount":      DeleteAccount,
				"deleteBudget":       DeleteBudget,
				"deleteGoal":         DeleteGoal,
				"deleteTransaction":  DeleteTransaction,
				"deleteUser":         DeleteUser,
				"getAccounts":        GetAccounts,
				"getAccountIncome":   GetAccountIncome,
				"getAccountExpense":  GetAccountExpense,
				"getAccountBalance":  GetAccountBalance,
				"getBudgets":         GetBudgets,
				"getGoals":           GetGoals,
				"getTransactionById": GetTransactionById,
				"getTransactions":    GetTransactions,
				"getCharts":          getCharts,
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
				Send_Error(ws, "Unknown Action", err)
			}
		}

		defer func() {
			slog.Info("Connection Closed", slog.String("user", userID))
			ws.Close()
			socket.IsActive = false
			if err := db.DB.Save(socket).Error; err != nil {
				Send_Error(ws, "Failed to update", err)
			}
		}()
	})(c)
}
