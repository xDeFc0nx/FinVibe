package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log/slog"
	"time"

	"github.com/pkg/errors"
	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func CreateWebSocket(userID string) (string, error) {
	socket := new(types.WebSocket)
	socket.ConnectionID = uuid.New().String()
	socket.IsActive = true
	socket.LastPing = time.Now()
	socket.CreatedAt = time.Now()

	var userExists bool
	checkErr := db.DB.QueryRow(context.Background(), `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`, userID).Scan(&userExists)
	if checkErr != nil {

		slog.Error(
			"Failed to check user existence for websocket",
			slog.String("userID", userID),
			slog.String("stack", fmt.Sprintf("%+v", errors.Wrap(checkErr, ""))))
		return "", nil
	}
	if !userExists {
		slog.Error("user with ID %s not found for webscoket", "UserID", userID)
		return "", nil
	}
	if _, err := db.DB.Exec(context.Background(), `
		INSERT INTO web_sockets (id, user_id, connection_id, is_active, last_ping, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)

			`, userID, userID, socket.ConnectionID, socket.IsActive, socket.LastPing, socket.CreatedAt); err != nil {
		slog.Error(
			"Error creating WebSocket connection",
			slog.String("userID", userID),
			slog.String("stack", fmt.Sprintf("%+v", errors.Wrap(err, ""))))
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
		data := map[string]any{
			"message": MsgUserNotFound,
		}
		SendError(ws, data, err)
		return

	}
	if _, err := db.DB.Exec(context.Background(), `
		UPDATE web_sockets
		SET
    last_ping = $1
		WHERE user_id = $2
		`, time.Now().UTC(), userID); err != nil {
		SendError(ws, MsgWebSocketUpdateFailed, err)
		return

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
			JSendFail(c, MsgWebSocketUpdateFailed, fiber.StatusBadRequest, err)
			return nil
		}

		return c.Status(101).
			JSON(fiber.Map{"message": "WebSocket connection updated successfully"})

	}
	data := map[string]any{
		"message": MsgUnauthorized,
	}
	JSendError(c, data, fiber.StatusUnauthorized, nil)
	return nil
}

func HandleWebSocketConnection(c *fiber.Ctx) error {
	slog.Info("Incoming Socket connection", slog.String("IP Address", c.IP()))

	socket := new(types.WebSocket)
	token := c.Cookies("jwt-token")

	if token == "" {
		data := map[string]any{
			"message": MsgMissingToken,
		}
		JSendError(c, data, fiber.StatusUnauthorized, nil)
		return nil
	}

	userID, connectionID, err := DecodeJWTToken(token)
	if err != nil {
		data := map[string]any{
			"message": MsgTokenDecodeFailed,
		}
		JSendError(c, data, fiber.StatusUnauthorized, nil)
		return nil
	}
	if err = db.DB.QueryRow(
		context.Background(), `
		SELECT id, is_active, last_ping
		FROM web_sockets
		WHERE connection_id = $1
		`,
		connectionID,
	).Scan(&socket.ID, &socket.IsActive, &socket.LastPing); err != nil {
		data := map[string]any{
			"message": fmt.Sprintf(MsgFetchFailedFmt, "websocket"),
		}
		JSendError(c, data, fiber.StatusInternalServerError, err)
		return err

	}
	if _, err := db.DB.Exec(context.Background(), `
		UPDATE web_sockets 
		SET
		is_active = true,
		last_ping = $1
		WHERE user_id = $2
		`, time.Now().UTC(), userID); err != nil {
		JSendFail(c, MsgWebSocketUpdateFailed, fiber.StatusBadRequest, err)
		return nil
	}
	return websocket.New(func(ws *websocket.Conn) {
		slog.Info(
			"User Connected",
			slog.String("UserID", userID),
		)
		done := make(chan struct{})
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
							SendMessage(ws, string(responseJson))
						}
					}

				case <-timeout:
					if socket.IsActive {
						if time.Since(socket.LastPing) > 20*time.Second {
							response := map[string]string{
								"message": "Timeout",
							}
							responseJson, _ := json.Marshal(response)
							SendMessage(ws, string(responseJson))
							ws.Close()
							return
						}
					}
				case <-done:
					return
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
				SendError(ws, MsgInvalidData, err)
				continue
			}
			slog.Info(
				"recv",
				slog.String("action", message.Action),
				"data", message.Data,
				slog.String("from", userID),
			)

			if message.Action == "" {
				SendError(ws, "Action required", err)
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
				"getAccount": 				GetAccount,
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
				SendError(ws, "Unknown Action", err)
			}
		}

		defer func() {
			slog.Info("Connection Closed", slog.String("user", userID))
			if _, err := db.DB.Exec(context.Background(), `
		UPDATE web_sockets 
		SET
		is_active = false
		WHERE user_id = $1
		`, userID); err != nil {
				data := map[string]any{
					"message": MsgWebSocketUpdateFailed,
				}
				SendError(ws, data, err)
				return
			}
			close(done)
			ws.Close()

		}()
	})(c)
}
