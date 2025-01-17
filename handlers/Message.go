package handlers

import (
	"encoding/json"
	"log/slog"

	"github.com/gofiber/contrib/websocket"
)

var InvalidData = "Error: Invalid form data"

func Send_Message(ws *websocket.Conn, sendText string) {
	response := map[string]string{
		"Message": sendText,
	}
	responseJSON, _ := json.Marshal(response)
	_ = ws.WriteMessage(websocket.TextMessage, responseJSON)
}

func Send_Error(ws *websocket.Conn, sendText string, err error) {
	response := map[string]interface{}{
		"Message": sendText,
	}

	if err != nil {
		response["Error"] = err.Error()
	}

	responseJSON, _ := json.Marshal(response)

	_ = ws.WriteMessage(websocket.TextMessage, responseJSON)

	if err != nil {
		slog.Error(sendText, slog.String("error", err.Error()))
	}
}
