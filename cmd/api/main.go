package main

import (
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

	if socket.IsActive == false {
		fmt.Printf("WebSocket is not active: ")
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	return websocket.New(func(ws *websocket.Conn) {
		for {
			mt, msg, err := ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			err = ws.WriteMessage(mt, msg)
			if err != nil {
				log.Println("write:", err)
				break
			}
			handlers.CreateTransaction()
		}

		defer func() {
			socket.IsActive = false
			socket.LastPing = time.Now()
			db.DB.Save(socket)
			ws.Close()
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
	app.Post("/Login", handlers.Login_func)
	app.Post("/logout", handlers.Logout_func)
	app.Get("/ws", HandleWebSocketConnection)

	app.Listen(":3000")

}
