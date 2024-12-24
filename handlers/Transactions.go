package handlers

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func CreateTransaction(conn *websocket.Conn, data json.RawMessage, userID string) {

	transaction := new(types.Transaction)

	transaction.ID = uuid.NewString()

	if err := json.Unmarshal(data, &transaction); err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"Invalid transaction data"}`))
		return
	}

	if err := db.DB.Create(&transaction).Error; err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"Failed to create transaction"}`))
		return
	}

	response, _ := json.Marshal(transaction)
	conn.WriteMessage(websocket.TextMessage, response)
}

func GetTransactions(conn *websocket.Conn, userID string) {
	transactions := []types.Transaction{}

	if err := db.DB.Where("user_id = ?", userID).Find(&transactions); err != nil {

		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"Failed to get transactions"}`))
	}

	response, _ := json.Marshal(transactions)
	conn.WriteMessage(websocket.TextMessage, response)
}
