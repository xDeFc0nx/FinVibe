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

func CreateRecurring(
	ws *websocket.Conn,
	data json.RawMessage,
	transactionID string,
	accountID string,
	userID string,
	amount float64,
) error {
	recurring := new(types.Recurring)
	recurring.ID = uuid.NewString()
	recurring.TransactionID = transactionID
	recurring.Amount = amount
	logger.Info("Raw data received: %s", string(data))

	if recurring.Frequency == "" {
		return ws.WriteMessage(
			websocket.TextMessage,
			[]byte(`{"Error":"Frequency is required"}`),
		)
	}
	if err := json.Unmarshal(data, &recurring); err != nil {
		Message(ws, "Error: Invalid form data")
		return err
	}

	logger.Info("Unmarshalled Recurring: %+v", recurring)

	recurring.StartDate = time.Now().Truncate(24 * time.Hour)

	switch recurring.Frequency {
	case "Daily":
		recurring.NextDate = time.Date(
			recurring.StartDate.Year(),
			recurring.StartDate.Month(),
			recurring.StartDate.Day()+1,
			0,
			0,
			0,
			0,
			recurring.StartDate.Location(),
		)

	case "Weekly":

		daysUntilNextWeek := int(time.Monday - recurring.StartDate.Weekday())
		if daysUntilNextWeek <= 0 {
			daysUntilNextWeek += 7
		}
		recurring.NextDate = recurring.StartDate.AddDate(
			0,
			0,
			daysUntilNextWeek,
		)

	case "Monthly":
		recurring.NextDate = time.Date(
			recurring.StartDate.Year(),
			recurring.StartDate.Month()+1,
			1,
			0,
			0,
			0,
			0,
			recurring.StartDate.Location(),
		)

	default:
		return ws.WriteMessage(
			websocket.TextMessage,
			[]byte(`{"Error":"Invalid frequency"}`),
		)
	}

	if err := db.DB.Save(&recurring).Error; err != nil {
		Message(ws, "Error: failed to Save recurring transaction")
		return err
	}
	go func() {
		err := handleRecurringTransaction(ws, recurring, userID, accountID)
		if err != nil {
			logger.Error(
				"Failed to handle recurring transaction: %s",
				err.Error(),
			)
		}
	}()

	return nil
}

func CreateTransaction(
	ws *websocket.Conn,
	data json.RawMessage,
	userID string,
) {
	transaction := new(types.Transaction)
	account := new(types.Accounts)
	recurring := new(types.Recurring)
	transaction.ID = uuid.NewString()

	transaction.UserID = userID

	if err := json.Unmarshal(data, &transaction); err != nil {
		Message(ws, "Error: Invalid form data")
		return
	}

	if transaction.AccountID == "" {
		Message(ws, "Error: Account ID is required")
		return
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, transaction.AccountID).First(&account).Error; err != nil {
		Message(ws, "Error: AccountID not found")
		return
	}
	if err := db.DB.Create(&transaction).Error; err != nil {
		Message(ws, "Error: Failed to create transaction")
		return
	}

	if transaction.IsRecurring {
		recurring.ID = uuid.NewString()
		recurring.TransactionID = transaction.ID
		var inputData map[string]interface{}
		if err := json.Unmarshal(data, &inputData); err != nil {

			Message(ws, "Error: Failed to parse recurring frequency")
			return
		}

		if freq, ok := inputData["Frequency"].(string); ok && freq != "" {
			recurring.Frequency = freq
		} else {

			Message(ws, "Error: Recurring Frequency is required")
			return
		}

		if err := db.DB.Create(recurring).Error; err != nil {
			Message(ws, "Error: Failed to create recurring transaction")
		}
		if err := CreateRecurring(ws, data, transaction.ID, account.ID, userID, transaction.Amount); err != nil {
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

	if err := ws.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}
}

func handleRecurringTransaction(
	ws *websocket.Conn,
	recurring *types.Recurring,
	userID string,
	accountID string,
) error {
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
			Message(ws, "Error: Failed to create new recurring transaction")

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
			Message(ws, "Error: Failed to update recurring next date")
			return nil
		}
	}
}

func GetTransactions(ws *websocket.Conn, data json.RawMessage, userID string) {
	account := new(types.Accounts)

	// Assuming the data contains a single object (AccountID)
	var requestData struct {
		AccountID string `json:"AccountID"`
	}

	if err := json.Unmarshal(data, &requestData); err != nil {
		Message(ws, "Error: Invalid form data")
		return
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, requestData.AccountID).First(&account).Error; err != nil {
		Message(ws, "Error: Account not found")
		return
	}

	transactions := []types.Transaction{}
	if err := db.DB.Where("account_id = ?", requestData.AccountID).Find(&transactions).Error; err != nil {
		Message(ws, "Error: Could Not get transactions")
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
	if err := ws.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}
}

func GetTransactionById(
	ws *websocket.Conn,
	data json.RawMessage,
	userID string,
) {
	transaction := new(types.Transaction)

	transaction.UserID = userID

	if err := json.Unmarshal(data, &transaction); err != nil {
		Message(ws, "Error: Invalid form data")
		return
	}

	if transaction.UserID != userID {
		Message(ws, "Error: Transaction does not belong to you")
		return
	}

	if err := db.DB.Where("id = ?", transaction.ID).First(&transaction).Error; err != nil {
		Message(ws, "Error: Transaction not found")
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
	if err := ws.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}
}

func UpdateTransaction(
	ws *websocket.Conn,
	data json.RawMessage,
	userID string,
) {
	transaction := new(types.Transaction)

	if err := json.Unmarshal(data, &transaction); err != nil {
		if err := ws.WriteMessage(websocket.TextMessage, []byte(`{"Error": "Invalid Data"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
			Message(ws, "Error: Invalid form data")
		}
	}

	if transaction.UserID != userID {
		Message(ws, "Error: Transaction does not byte to you")
		return
	}

	if err := db.DB.Where("id = ?", transaction.ID).First(&transaction).Error; err != nil {
		Message(ws, "Error: Transaction not found")
		return

	}

	if err := db.DB.Save(transaction).Error; err != nil {
		Message(ws, "Error: Failed to save")
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
	if err := ws.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}
}

func DeleteTransaction(
	ws *websocket.Conn,
	data json.RawMessage,
	userID string,
) {
	transaction := new(types.Transaction)

	transaction.UserID = userID

	if err := json.Unmarshal(data, &transaction); err != nil {
		Message(ws, "Error: Invalid form data")

		return
	}
	if transaction.UserID != userID {
		Message(ws, "Error: Transaction does not belong to you")
		return
	}

	if err := db.DB.Where("id = ?", transaction.ID).First(&transaction).Error; err != nil {
		Message(ws, "Error: Transaction not found")
		return
	}

	if err := db.DB.Delete(transaction).Error; err != nil {
		Message(ws, "Error: Failed to delete")
		return
	}

	Message(ws, "Success: Transaction Deleted")
}

func GetAccountBalance(ws *websocket.Conn, accountID string) error {
	transactions := []types.Transaction{}
	account := new(types.Accounts)

	account.ID = accountID

	if err := db.DB.Where(" id =?", account.ID).First(&account).Error; err != nil {
		Message(ws, "Error: Account not found")
		return err
	}

	if err := db.DB.Where("account_id = ?", account.ID).Find(&transactions).Error; err != nil {
		Message(ws, "Error: Could not get transactions")
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
		Message(ws, "Error: Failed to save")
		return err
	}

	return nil
}
