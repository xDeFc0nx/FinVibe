package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/xDeFc0nx/logger-go-pkg"

	"github.com/xDeFc0nx/FinVibe/cmd/flag"
	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/handlers"
	"github.com/xDeFc0nx/FinVibe/types"
)

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

	if err := handlers.HandleCheckAuth(c, userID); err != nil {
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
				select {
				case <-ticker.C:

					if time.Since(socket.LastPing) > 15*time.Second {
						if err := c.WriteMessage(websocket.TextMessage, []byte(`ping`)); err != nil {
							logger.Debug("Error sending ping: %v", err)
							return
						}

					}

				case <-time.After(15 * time.Second):
					if time.Since(socket.LastPing) > 15*time.Second {
						if err := c.WriteMessage(websocket.TextMessage, []byte(`"error": "Connection timeout"`)); err != nil {
							logger.Debug("Error sending ping: %v", err)
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

			switch message.Action {
			case "pong":
				handlers.HeartBeat(c, userID)
			case "getUser":
				handlers.GetUser(c, userID)
			case "updateUser":
				handlers.UpdateUser(c, message.Data, userID)
			case "deleteUser":
				handlers.DeleteUser(c, userID)
			case "logout":
				handlers.LogoutHandler(c, userID)
			case "createAccount":
				handlers.CreateAccount(c, message.Data, userID)
			case "getAccounts":
				handlers.GetAccounts(c, userID)
			case "updateAccount":
				handlers.UpdateAccount(c, message.Data, userID)
			case "deleteAccount":
				handlers.DeleteAccount(c, message.Data, userID)
			case "createTransaction":
				handlers.CreateTransaction(c, message.Data, userID)
			case "getTransactions":
				handlers.GetTransactions(c, message.Data, userID)
			case "getTransactionById":
				handlers.GetTransactionById(c, message.Data, userID)
			case "updateTransaction":
				handlers.UpdateTransction(c, message.Data, userID)
			case "deleteTransaction":
				handlers.DeleteTransaction(c, message.Data, userID)
			case "createBudget":
				handlers.CreateBudget(c, message.Data, userID)
			case "getBudgets":
				handlers.GetBudgets(c, userID)
			case "updateBudget":
				handlers.UpdateBudget(c, message.Data, userID)
			case "deleteBudget":
				handlers.DeleteBudget(c, message.Data, userID)
			case "creteGoal":
				handlers.CreateGoal(c, message.Data, userID)
			case "getGoals":
				handlers.Getgoals(c, userID)
			case "updateGoals":
				handlers.Updategoal(c, message.Data, userID)
			case "deleteGoal":
				handlers.DeleteGoal(c, message.Data, userID)

			default:
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

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow any origin
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	db.Conn()

	flag.Flag()

	app.Post("/Register", handlers.CreateUser)
	app.Post("/Login", handlers.LoginHandler)
	app.Post("/checkAuth", handlers.CheckAuth)
	app.Use("/ws", HandleWebSocketConnection)

	if err := app.Listen(os.Getenv("PORT")); err != nil {
		logger.Error("%s", err.Error())
	}

}
