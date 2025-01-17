package handlers

import (
	"encoding/json"
	"log/slog"

	"github.com/gofiber/contrib/websocket"
)

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
