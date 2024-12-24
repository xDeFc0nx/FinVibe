package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func CreateTransaction(conn *websocket.Conn, data json.RawMessage, userID string) {

	transaction := new(types.Transaction)

	transaction.ID = uuid.NewString()

	transaction.UserID = userID

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

	if err := db.DB.Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		fmt.Printf("Database error: %v\n", err)
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"Failed to get transactions"}`))
		return
	}

	response, _ := json.Marshal(transactions)
	conn.WriteMessage(websocket.TextMessage, response)
}

func GetTransactionById(conn *websocket.Conn, data json.RawMessage, userID string) {
	transaction := new(types.Transaction)

	transaction.UserID = userID

	if err := json.Unmarshal(data, &transaction); err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"Invalid transaction data"}`))
		return
	}
	if transaction.UserID != userID {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"Transaction does not belong to the user"}`))
		return
	}

	if err := db.DB.Where("id = ?", transaction.ID).First(&transaction).Error; err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"Transaction not found"}`))
		return
	}

	response, _ := json.Marshal(transaction)
	conn.WriteMessage(websocket.TextMessage, response)
}
