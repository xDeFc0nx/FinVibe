package handlers

import (
	"encoding/json"
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"fmt"
	"github.com/pkg/errors"

	"github.com/gofiber/contrib/websocket"
)

type Error struct {
}

type Response struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

func JSendSuccess(c *fiber.Ctx, data any) {
	c.Status(fiber.StatusOK).JSON(Response{
		Status: "success",
		Data:   data,
	})

}

func JSendFail(c *fiber.Ctx, data any, code int, err error) {
	c.Status(code).JSON(Response{
		Status: "fail",
		Data:   data,
	})
	slog.Info("fail",
		data,
		slog.String("stack", fmt.Sprintf("%+v", errors.Wrap(err, ""))),
	)

}

func JSendError(c *fiber.Ctx, data any, code int, err error) {
	c.Status(code).JSON(Response{
		Status: "error",
		Data:   data,
	})

	slog.Error("error",
		data,
		slog.String("stack", fmt.Sprintf("%+v", errors.Wrap(err, ""))),
	)
}

var InvalidData = "Error: Invalid form data"

func SendMessage(ws *websocket.Conn, sendText string) {
	if err := ws.WriteMessage(websocket.TextMessage, []byte(sendText)); err != nil {
		slog.Error("failed to send message", slog.String("Error", err.Error()))
	}
}

func SendError(ws *websocket.Conn, sendText string, err error) {
	slog.Error("error",
		slog.String("stack", fmt.Sprintf("%+v", errors.Wrap(err, ""))),
	)

	response := map[string]any{
		"Error": sendText,
	}
	responseJSON, _ := json.Marshal(response)
	if err := ws.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
		slog.Error(sendText, fmt.Sprintf("%+v", errors.Wrap(err, "")))
	}
}
