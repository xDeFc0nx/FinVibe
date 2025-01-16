package handlers

import (
	"encoding/json"
	"log/slog"

	"github.com/gofiber/contrib/websocket"
)

var InvalidData = "Error: Invalid form data"

func Message(ws *websocket.Conn, sendText string, err error) {
	textJsonify, _ := json.Marshal(sendText)
	if err := ws.WriteMessage(websocket.TextMessage, []byte(textJsonify)); err != nil {
		slog.Error("failed to send message", slog.String("error", err.Error()))
	}
	if err != nil {
		slog.Error(sendText, slog.String("error", err.Error()))
	}
}
