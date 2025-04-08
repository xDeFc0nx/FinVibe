package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func CreateWebSocket(userID string) (string, error) {
	socket := new(types.WebSocket)
	socket.UserID = userID
	socket.ConnectionID = uuid.New().String()
	socket.IsActive = true
	socket.LastPing = time.Now()
	socket.CreatedAt = time.Now()

	var userExists bool
	checkErr := db.DB.QueryRow(context.Background(), `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`, userID).Scan(&userExists)
	if checkErr != nil {
		slog.Error("Failed to check user existence for websocket", slog.String("error", checkErr.Error()), slog.String("userID", userID))
		return "", nil
	}
	if !userExists {
		slog.Error("user with ID %s not found for webscoket", userID)
		return "", nil
	}
	if _, err := db.DB.Exec(context.Background(), `
		INSERT INTO web_sockets (id, user_id, connection_id, is_active, last_ping, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)

			`, userID, userID, socket.ConnectionID, socket.IsActive, socket.LastPing, socket.CreatedAt); err != nil {
		slog.Error(
			"Error creating WebSocket connection",
			"Err",
			err,
		)
	}

	return socket.ID, nil
}

func HeartBeat(ws *websocket.Conn, data json.RawMessage, userID string) {

	var userExists bool
	err := db.DB.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)",
		userID,
	).Scan(&userExists)
	if err != nil || !userExists {
		Send_Error(ws, "User not found", err)
		return

	}
	if _, err := db.DB.Exec(context.Background(), `
		UPDATE web_sockets 
    last_ping = $1
		WHERE user_id = $2

		`, userID, time.Now().UTC()); err != nil {
		Send_Error(ws, "Failed to update last ping", err)

	}
}
func HandleCheckAuth(c *fiber.Ctx, userID string) error {
	if CheckAuth(c) == nil {

		var userExists bool
		err := db.DB.QueryRow(
			context.Background(),
			"SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)",
			userID,
		).Scan(&userExists)
		if err != nil || !userExists {
			return err

		}
		if _, err := db.DB.Exec(context.Background(), `
		UPDATE web_sockets 
		is_active = true,
    last_ping = $1
		WHERE user_id = $2

		`, userID, time.Now().UTC()); err != nil {
			JSendFail(c, "Failed to update last ping", 500)
		}

		return c.Status(101).
			JSON(fiber.Map{"message": "WebSocket connection updated successfully"})
	}
	return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
}

func HandleWebSocketConnection(c *fiber.Ctx) error {
	slog.Info("Incoming Socket connection", slog.String("IP Address", c.IP()))

	socket := new(types.WebSocket)
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
	var userExists bool
	err = db.DB.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM web_sockets WHERE connection_id = $1)",
		connectionID,
	).Scan(&userExists)
	if err != nil || !userExists {
		return err

	}
	if _, err := db.DB.Exec(context.Background(), `
		UPDATE web_sockets 
		is_active = true,
    last_ping = $1
		WHERE user_id = $2

		`, time.Now().UTC(), userID); err != nil {
		JSendFail(c, "Failed to update last ping", 500)
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
				"data", message.Data,
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
			if _, err := db.DB.Exec(context.Background(), `
		UPDATE web_sockets 
		is_active = false,
		WHERE user_id = $1

		`, userID); err != nil {
				JSendFail(c, "Failed to update last ping", 500)
			}
		}()
	})(c)
}
