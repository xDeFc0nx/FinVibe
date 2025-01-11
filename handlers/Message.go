package handlers

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
	"github.com/xDeFc0nx/logger-go-pkg"
)

var InvalidData = "Error: Invalid form data"

func Message(ws *websocket.Conn, sendText string) {
	textJsonify, _ := json.Marshal(sendText)
	if err := ws.WriteMessage(websocket.TextMessage, []byte(textJsonify)); err != nil {
		logger.Error("%s", err.Error())
	}
}
