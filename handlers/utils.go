package handlers

import (
	"encoding/json"
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/contrib/websocket"
)

type Error struct{



}


type Response struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

func JSendSuccess(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Status: "success",
		Data:   data,
	})
}

func JSendFail(c *fiber.Ctx, data any, code int) error {
	return c.Status(code).JSON(Response{
		Status: "fail",
		Data:   data,
	})
}

func JSendError(c *fiber.Ctx, data any, code int) error {
	return c.Status(code).JSON(Response{
		Status: "error",
		Data:   data,
	})
}

var InvalidData = "Error: Invalid form data"

func Send_Message(ws *websocket.Conn, sendText string) {
	if err := ws.WriteMessage(websocket.TextMessage, []byte(sendText)); err != nil {
		slog.Error("failed to send message", slog.String("Error", err.Error()))
	}
}

func Send_Error(ws *websocket.Conn, sendText string, err error) {
	response := map[string]interface{}{
		"Error": sendText,
	}
	responseJSON, _ := json.Marshal(response)
	if err := ws.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
		slog.Error(sendText, slog.String("error", err.Error()))
	}
}
