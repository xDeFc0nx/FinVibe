package handlers

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/xDeFc0nx/logger-go-pkg"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func CreateTransaction(c *websocket.Conn, data json.RawMessage, userID string) {

	transaction := new(types.Transaction)
	account := new(types.Accounts)
	transaction.ID = uuid.NewString()

	transaction.UserID = userID
	transaction.AccountID = account.ID

	if err := json.Unmarshal(data, &transaction); err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid transaction data"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	if transaction.AccountID == "" {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Account ID is required"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, transaction.AccountID).First(&account).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Account not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}
	if err := db.DB.Create(&transaction).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to create transaction"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	response := map[string]interface{}{
		"Success": "Created Transaction",
		"Transaction": map[string]string{
			"ID":        transaction.ID,
			"UserID":    transaction.UserID,
			"AccountID": transaction.AccountID,
		},
	}

	responseData, _ := json.Marshal(response)

	if err := c.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}

}
func GetTransactions(c *websocket.Conn, data json.RawMessage, userID string) {
	account := new(types.Accounts)

	// Assuming the data contains a single object (AccountID)
	var requestData struct {
		AccountID string `json:"AccountID"`
	}

	if err := json.Unmarshal(data, &requestData); err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid request data"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, requestData.AccountID).First(&account).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Account not found"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	transactions := []types.Transaction{}
	if err := db.DB.Where("account_id = ?", requestData.AccountID).Find(&transactions).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Could not fetch transactions"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	transactionData := make([]map[string]interface{}, len(transactions))

	for i, t := range transactions {
		transactionData[i] = map[string]interface{}{
			"ID":        t.ID,
			"UserID":    t.UserID,
			"AccountID": t.AccountID,
			"Amount":    t.Amount,
		}
	}

	// Package the response
	response := map[string]interface{}{
		"Success":      "Fetched Transactions",
		"Transactions": transactionData,
	}

	// Marshal the response and send it
	responseData, _ := json.Marshal(response)
	if err := c.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}

}

func GetTransactionById(c *websocket.Conn, data json.RawMessage, userID string) {
	transaction := new(types.Transaction)

	transaction.UserID = userID

	if err := json.Unmarshal(data, &transaction); err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid transaction data"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}
	if transaction.UserID != userID {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction does not belong to the user"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	if err := db.DB.Where("id = ?", transaction.ID).First(&transaction).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction not found"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	response, _ := json.Marshal(transaction)
	if err := c.WriteMessage(websocket.TextMessage, response); err != nil {
		logger.Error("%s", err.Error())
	}
}

func UpdateTransction(c *websocket.Conn, data json.RawMessage, userID string) {
	transaction := new(types.Transaction)

	if err := json.Unmarshal(data, &transaction); err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error": "Invalid Data"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
	}

	if transaction.UserID != userID {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction does not belong to the user"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	if err := db.DB.Where("id = ?", transaction.ID).First(&transaction).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return

	}

	if err := db.DB.Save(transaction).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to Update"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
	}
	if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Success":"Transaction Updated"}`)); err != nil {
		logger.Error("%s", err.Error())
	}
}
func DeleteTransaction(c *websocket.Conn, data json.RawMessage, userID string) {
	transaction := new(types.Transaction)

	transaction.UserID = userID

	if err := json.Unmarshal(data, &transaction); err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid transaction data"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

		return
	}
	if transaction.UserID != userID {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction does not belong to the user"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	if err := db.DB.Where("id = ?", transaction.ID).First(&transaction).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	if err := db.DB.Delete(transaction).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"failed to Delete"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Success":"Transaction Deleted"}`)); err != nil {
		logger.Error("%s", err.Error())
	}

}

func GetAccountBalance(c *websocket.Conn, accountID string) error {
	transactions := []types.Transaction{}
	account := new(types.Accounts)

	account.ID = accountID

	if err := db.DB.Where(" id =?", account.ID).First(&account).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Account not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return err
	}

	if err := db.DB.Where("account_id = ?", account.ID).Find(&transactions).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Could not fetch transactions"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
		return err
	}

	totalBalance := float64(0)
	for _, t := range transactions {
		if t.AccountID == accountID {
			totalBalance += t.Amount
		}
	}

	account.Balance = totalBalance

	if err := db.DB.Save(account).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"failed to Save"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return err
	}

	return nil
}
