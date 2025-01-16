package main

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/xDeFc0nx/FinVibe/cmd/flag"
	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/handlers"
)

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))
	db.Conn()

	flag.Flag()

	app.Post("/Register", handlers.CreateUser)
	app.Post("/Login", handlers.LoginHandler)
	app.Post("/logout", handlers.LogoutHandler)
	app.Get("/checkAuth", handlers.CheckAuth)
	app.Use("/ws", handlers.HandleWebSocketConnection)

	app.Static("/static", "./client/dist/static", fiber.Static{
		Compress: true,
	})
	app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendFile("./client/dist/index.html")
	})
	if err := app.Listen(os.Getenv("PORT")); err != nil {
		slog.Error("Error Listening", slog.String("error", err.Error()))
	}
}
