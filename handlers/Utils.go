package handlers

import (
	"encoding/json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/rotisserie/eris"
	"log/slog"
)

func Eris(err error, status string, data any) {

	wrapError := eris.Wrap(err, "error occurred during operation")
	errorJson := eris.ToJSON(wrapError, true)
	slog.Error(status,
		slog.Any("message", data),
		slog.Any(status, errorJson),
	)
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
	Eris(err, "fail", data)
	c.Status(code).JSON(Response{
		Status: "fail",
		Data:   data,
	})
	Eris(err, "fail", data)
}

func JSendError(c *fiber.Ctx, data any, code int, err error) {
	c.Status(code).JSON(Response{
		Status: "error",
		Data:   data,
	})

	Eris(err, "error", data)

}

func SendMessage(ws *websocket.Conn, sendText string) {
	if err := ws.WriteMessage(websocket.TextMessage, []byte(sendText)); err != nil {
		slog.Error("failed to send message", slog.String("Error", err.Error()))
	}
}

func SendError(ws *websocket.Conn, data any, err error) {
	Eris(err, "error", data)
	response := map[string]any{
		"Status": "error",
		"Data":   data,
	}
	responseJSON, _ := json.Marshal(response)
	if err := ws.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
		slog.Error("failed to send error")

	}
}
