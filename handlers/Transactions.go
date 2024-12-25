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
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid transaction data"}`))
		return
	}

	if err := db.DB.Create(&transaction).Error; err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to create transaction"}`))
		return
	}

	response, _ := json.Marshal(transaction)
	conn.WriteMessage(websocket.TextMessage, response)
}

func GetTransactions(conn *websocket.Conn, userID string) {
	transactions := []types.Transaction{}

	if err := db.DB.Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		fmt.Printf("Database Error: %v\n", err)
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to get transactions"}`))
		return
	}

	response, _ := json.Marshal(transactions)
	conn.WriteMessage(websocket.TextMessage, response)
}

func GetTransactionById(conn *websocket.Conn, data json.RawMessage, userID string) {
	transaction := new(types.Transaction)

	transaction.UserID = userID

	if err := json.Unmarshal(data, &transaction); err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid transaction data"}`))
		return
	}
	if transaction.UserID != userID {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction does not belong to the user"}`))
		return
	}

	if err := db.DB.Where("id = ?", transaction.ID).First(&transaction).Error; err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction not found"}`))
		return
	}

	response, _ := json.Marshal(transaction)
	conn.WriteMessage(websocket.TextMessage, response)
}

func UpdateTransction(conn *websocket.Conn, data json.RawMessage, userID string) {
	transaction := new(types.Transaction)

	if err := json.Unmarshal(data, &transaction); err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error": "Invalid Data"}`))
	}

	if transaction.UserID != userID {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction does not belong to the user"}`))
		return
	}

	if err := db.DB.Where("id = ?", transaction.ID).First(&transaction).Error; err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction not found"}`))
		return
	}

	if err := db.DB.Save(transaction).Error; err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to Update"}`))
	}
	conn.WriteMessage(websocket.TextMessage, []byte(`{"Success":"Transaction Updated"}`))
}
func DeleteTransaction(conn *websocket.Conn, data json.RawMessage, userID string) {
	transaction := new(types.Transaction)

	transaction.UserID = userID

	if err := json.Unmarshal(data, &transaction); err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid transaction data"}`))
		return
	}
	if transaction.UserID != userID {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction does not belong to the user"}`))
		return
	}

	if err := db.DB.Where("id = ?", transaction.ID).First(&transaction).Error; err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction not found"}`))
		return
	}

	if err := db.DB.Delete(transaction).Error; err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"failed to Delete"}`))
		return
	}

	conn.WriteMessage(websocket.TextMessage, []byte(`{"Success":"Transaction Deleted"}`))

}
