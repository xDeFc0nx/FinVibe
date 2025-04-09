package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/handlers"
	"io"
	"log/slog"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading env", slog.String("err", err.Error()))
	}
	file, err := os.OpenFile(
		"app.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		slog.Error("Error opening log file", slog.String("error", err.Error()))
		return
	}
	defer file.Close()

	multiWriter := io.MultiWriter(file, os.Stdout)

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(multiWriter, opts)

	logger := slog.New(handler)

	slog.SetDefault(logger)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))
	db.Conn()

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
		slog.Error("Error Listening",
			slog.String("stack", fmt.Sprintf("%+v", errors.Wrap(err, ""))))
	}
}
