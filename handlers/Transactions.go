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
	account := new(types.Accounts)
	transaction.ID = uuid.NewString()

	transaction.UserID = userID

	if err := json.Unmarshal(data, &transaction); err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid transaction data"}`))
		return
	}

	if err := db.DB.Where("user_id =? AND account_id =?", userID, account.ID).First(&account).Error; err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Account not found"}`))
		return
	}
	if err := db.DB.Create(&transaction).Error; err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to create transaction"}`))
		return
	}

	response, _ := json.Marshal(transaction)
	conn.WriteMessage(websocket.TextMessage, response)
}

func GetTransactions(conn *websocket.Conn, data json.RawMessage, userID string) {
	transactions := []types.Transaction{}
	account := new(types.Accounts)

	if err := json.Unmarshal(data, &transactions); err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid transaction data"}`))
		return
	}
	if err := db.DB.Where("user_id =? AND account_id =?", userID, account.ID).First(&account).Error; err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Account not found"}`))
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
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error": "Invalid Data"}`+err.Error()))
	}

	if transaction.UserID != userID {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction does not belong to the user"}`))
		return
	}

	if err := db.DB.Where("id = ?", transaction.ID).First(&transaction).Error; err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction not found"}`+err.Error()))
		return
	}

	if err := db.DB.Save(transaction).Error; err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to Update"}`+err.Error()))
	}
	conn.WriteMessage(websocket.TextMessage, []byte(`{"Success":"Transaction Updated"}`))
}
func DeleteTransaction(conn *websocket.Conn, data json.RawMessage, userID string) {
	transaction := new(types.Transaction)

	transaction.UserID = userID

	if err := json.Unmarshal(data, &transaction); err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid transaction data"}`+err.Error()))
		return
	}
	if transaction.UserID != userID {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction does not belong to the user"}`))
		return
	}

	if err := db.DB.Where("id = ?", transaction.ID).First(&transaction).Error; err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction not found"}`+err.Error()))
		return
	}

	if err := db.DB.Delete(transaction).Error; err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"Error":"failed to Delete"}`+err.Error()))
		return
	}

	conn.WriteMessage(websocket.TextMessage, []byte(`{"Success":"Transaction Deleted"}`))

}
