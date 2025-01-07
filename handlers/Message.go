package handlers

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/xDeFc0nx/logger-go-pkg"
)

func Message(ws *websocket.Conn, sendText string) {
	if err := ws.WriteMessage(websocket.TextMessage, []byte(sendText)); err != nil {
		logger.Error("%s", err.Error())
	}
}
