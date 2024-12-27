package handlers

import (
	"encoding/json"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/xDeFc0nx/logger-go-pkg"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func CreateRecurring(c *websocket.Conn, data json.RawMessage, transactionID string, accountID string, userID string, amount float64) error {
	recurring := new(types.Recurring)
	recurring.ID = uuid.NewString()
	recurring.TransactionID = transactionID
	recurring.Amount = amount
	logger.Info("Raw data received: %s", string(data))

	if recurring.Frequency == "" {
		return c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Frequency is required"}`))
	}
	if err := json.Unmarshal(data, &recurring); err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid recurring data"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return err
	}

	logger.Info("Unmarshalled Recurring: %+v", recurring)

	recurring.StartDate = time.Now().Truncate(24 * time.Hour)

	switch recurring.Frequency {
	case "Daily":
		recurring.NextDate = time.Date(recurring.StartDate.Year(), recurring.StartDate.Month(), recurring.StartDate.Day()+1, 0, 0, 0, 0, recurring.StartDate.Location())

	case "Weekly":

		daysUntilNextWeek := int(time.Monday - recurring.StartDate.Weekday())
		if daysUntilNextWeek <= 0 {
			daysUntilNextWeek += 7
		}
		recurring.NextDate = recurring.StartDate.AddDate(0, 0, daysUntilNextWeek)

	case "Monthly":
		recurring.NextDate = time.Date(recurring.StartDate.Year(), recurring.StartDate.Month()+1, 1, 0, 0, 0, 0, recurring.StartDate.Location())

	default:
		return c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid frequency"}`))
	}

	if err := db.DB.Save(&recurring).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to save recurring transaction"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return err
	}
	go func() {
		err := handleRecurringTransaction(c, recurring, userID, accountID)
		if err != nil {
			logger.Error("Failed to handle recurring transaction: %s", err.Error())
		}
	}()

	return nil
}
func CreateTransaction(c *websocket.Conn, data json.RawMessage, userID string) {
	transaction := new(types.Transaction)
	account := new(types.Accounts)
	recurring := new(types.Recurring)
	transaction.ID = uuid.NewString()

	transaction.UserID = userID

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

	if transaction.IsRecurring {
		recurring.ID = uuid.NewString()
		recurring.TransactionID = transaction.ID
		var inputData map[string]interface{}
		if err := json.Unmarshal(data, &inputData); err != nil {
			if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to parse recurring frequency"}`+err.Error())); err != nil {
				logger.Error("%s", err.Error())
			}

			return
		}

		if freq, ok := inputData["Frequency"].(string); ok && freq != "" {
			recurring.Frequency = freq
		} else {

			if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Recurring Frequency is required for recurring transactions"}`)); err != nil {
				logger.Error("%s", err.Error())
			}
			return
		}

		if err := db.DB.Create(recurring).Error; err != nil {
			if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to create recurring transaction"}`+err.Error())); err != nil {
				logger.Error("%s", err.Error())
			}
			return
		}
		if err := CreateRecurring(c, data, transaction.ID, account.ID, userID, transaction.Amount); err != nil {
			return
		}
	}

	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID":          transaction.ID,
			"UserID":      transaction.UserID,
			"AccountID":   transaction.AccountID,
			"Amount":      transaction.Amount,
			"IsRecurring": transaction.IsRecurring,
			"Frequency":   recurring.Frequency,
		},
	}

	responseData, _ := json.Marshal(response)

	if err := c.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}

}
func handleRecurringTransaction(c *websocket.Conn, recurring *types.Recurring, userID string, accountID string) error {
	for {
		time.Sleep(time.Until(recurring.NextDate))

		newTransaction := types.Transaction{
			ID:        uuid.NewString(),
			UserID:    userID,
			AccountID: accountID,
			Amount:    recurring.Amount,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.DB.Create(&newTransaction).Error; err != nil {
			if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to create new recurring transaction"}`+err.Error())); err != nil {
				logger.Error("%s", err.Error())
			}
			return nil
		}

		switch recurring.Frequency {
		case "Daily":
			recurring.NextDate = recurring.NextDate.Add(24 * time.Hour)
		case "Weekly":
			recurring.NextDate = recurring.NextDate.Add(7 * 24 * time.Hour)
		case "Monthly":
			recurring.NextDate = recurring.NextDate.AddDate(0, 1, 0)
		}

		if err := db.DB.Save(&recurring).Error; err != nil {
			if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to update recurring next date"}`+err.Error())); err != nil {
				logger.Error("%s", err.Error())
			}
			return nil
		}
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
			"ID":          t.ID,
			"UserID":      t.UserID,
			"AccountID":   t.AccountID,
			"Amount":      t.Amount,
			"IsRecurring": t.IsRecurring,
		}
	}

	// Package the response
	response := map[string]interface{}{
		"Success": transactionData,
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

	transactionData := map[string]interface{}{
		"ID":          transaction.ID,
		"UserID":      transaction.UserID,
		"AccountID":   transaction.AccountID,
		"Amount":      transaction.Amount,
		"IsRecurring": transaction.IsRecurring,
	}

	response := map[string]interface{}{
		"Success": transactionData,
	}

	responseData, _ := json.Marshal(response)
	if err := c.WriteMessage(websocket.TextMessage, responseData); err != nil {
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

	transactionData := map[string]interface{}{
		"ID":        transaction.ID,
		"UserID":    transaction.UserID,
		"AccountID": transaction.AccountID,
		"Amount":    transaction.Amount,
	}

	response := map[string]interface{}{
		"Success": transactionData,
	}

	responseData, _ := json.Marshal(response)
	if err := c.WriteMessage(websocket.TextMessage, responseData); err != nil {
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
