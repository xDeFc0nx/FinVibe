package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/xDeFc0nx/FinVibe/cmd/flag"
	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/handlers"
	"github.com/xDeFc0nx/FinVibe/types"
)

func HandleWebSocketConnection(c *fiber.Ctx) error {
	connectionID := c.Query("ConnectionID")
	socket := new(types.WebSocketConnection)

	if connectionID == "" {
		fmt.Printf("Connection ID is required")
		return c.Status(400).JSON(fiber.Map{"error": "ConnectionID is required"})

	}

	if err := db.DB.Where("id = ?", connectionID).First(socket).Error; err != nil {

		fmt.Printf("WebSocket not found: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "WebSocket not found"})

	}

	userID := socket.UserID

	token := c.Cookies("jwt-token")
	if token == "" {
		handlers.DecodeJWTToken(token)
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if socket.IsActive == false {
		fmt.Printf("WebSocket is not active")
		return c.Status(403).JSON(fiber.Map{"error": "WebSocket is not active"})
	}
	return websocket.New(func(c *websocket.Conn) {

		var (
			msg []byte
			err error
		)
		for {
			if _, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			// Parse the JSON message
			var message struct {
				Action string          `json:"action"`
				Data   json.RawMessage `json:"data"`
			}
			if err := json.Unmarshal(msg, &message); err != nil {
				c.WriteMessage(websocket.TextMessage, []byte(`{"error":"Invalid message format"}`))
				continue
			}

			switch message.Action {
			case "createTransaction":
				handlers.CreateTransaction(c, message.Data, userID)
			case "getUser":
				handlers.GetUser(c, userID)
			case "updateUser":
				handlers.UpdateUser(c, message.Data, userID)

			case "logout":
				handlers.LogoutHandler(c, userID)
			case "getTransactions":
				handlers.GetTransactions(c, userID)
			case "getTransactionsById":
				handlers.GetTransactionById(c, message.Data, userID)

			default:
				c.WriteMessage(websocket.TextMessage, []byte(`{"error":"Unknown action"}`))
			}
		}

		defer func() {
			socket.IsActive = false
			socket.LastPing = time.Now()
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

	app.Listen(":3000")

}
