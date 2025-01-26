package handlers

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

var requestData struct {
	AccountID string `json:"AccountID"`
	DateRange string `json:"DateRange"`
}

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

	if recurring.Frequency == "" {
		Send_Error(ws, "Invalid frequency", nil)
		return nil
	}
	if err := json.Unmarshal(data, &recurring); err != nil {
		Send_Error(ws, InvalidData, err)
		return err
	}

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
		Send_Error(ws, "Invalid frequency", nil)
		return nil
	}

	if err := db.DB.Save(&recurring).Error; err != nil {
		Send_Error(ws, "Error: failed to Save recurring transaction", err)
		return err
	}
	go func() {
		err := handleRecurringTransaction(ws, recurring, userID, accountID)
		if err != nil {
			slog.Error(
				"Failed to handle recurring transaction",
				slog.String("error", err.Error()),
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
	recurring := new(types.Recurring)
	transaction.ID = uuid.NewString()
	account := new(types.Accounts)

	transaction.UserID = userID
	if err := json.Unmarshal(data, &transaction); err != nil {
		Send_Error(ws, InvalidData, err)
		return
	}

	if transaction.AccountID == "" {
		Send_Error(ws, "Account ID is required", nil)
		return
	}

	if err := db.DB.Where("user_id = ? AND id = ?", userID, transaction.AccountID).First(&account).Error; err != nil {
		Send_Error(ws, "AccountID not found", err)
		return
	}
	if err := db.DB.Create(&transaction).Error; err != nil {
		Send_Error(ws, "Failed to create transaction", err)
		return
	}

	if transaction.IsRecurring {
		recurring.ID = uuid.NewString()
		recurring.TransactionID = transaction.ID

		var inputData map[string]interface{}
		if err := json.Unmarshal(data, &inputData); err != nil {

			Send_Error(ws, "Failed to parse recurring frequency", err)
			return
		}

		if freq, ok := inputData["Frequency"].(string); ok && freq != "" {
			recurring.Frequency = freq
		} else {

			Send_Error(ws, "Recurring Frequency is required", nil)
			return
		}

		if err := db.DB.Create(recurring).Error; err != nil {
			Send_Error(ws, "Failed to create recurring transaction", err)
		}
		if err := CreateRecurring(ws, data, transaction.ID, account.ID, userID, transaction.Amount); err != nil {
			return
		}
	}

	response := map[string]interface{}{
		"transaction": map[string]interface{}{
			"ID":          transaction.ID,
			"UserID":      transaction.UserID,
			"AccountID":   transaction.AccountID,
			"Type":        transaction.Type,
			"Amount":      transaction.Amount,
			"Description": transaction.Description,
			"IsRecurring": transaction.IsRecurring,
			"Frequency":   recurring.Frequency,
			"CreatedAt":   recurring.CreatedAt.Format(time.RFC3339),
		},
	}

	responseData, _ := json.Marshal(response)

	Send_Message(ws, string(responseData))
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
			Send_Error(ws, "Failed to create new recurring transaction", err)

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
			Send_Error(ws, "Failed to update recurring next date", err)
			return nil
		}
	}
}

func GetTransactions(ws *websocket.Conn, data json.RawMessage, userID string) {
	account := new(types.Accounts)

	if err := json.Unmarshal(data, &requestData); err != nil {
		Send_Error(ws, InvalidData, err)
		return
	}

	if requestData.DateRange == "" {
		Send_Error(ws, "Date Range is Required", nil)
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, requestData.AccountID).First(&account).Error; err != nil {
		Send_Error(ws, "Account not found", err)
		return
	}
	start, end := GetDateRange(requestData.DateRange)

	transactions := []types.Transaction{}
	if err := db.DB.Where("account_id = ? AND created_at BETWEEN ? AND ?", requestData.AccountID, start, end).Find(&transactions).Error; err != nil {
		Send_Error(ws, "Could Not get transactions", err)
		return
	}

	transactionData := make([]map[string]interface{}, len(transactions))

	for i, t := range transactions {
		transactionData[i] = map[string]interface{}{
			"ID":          t.ID,
			"UserID":      t.UserID,
			"AccountID":   t.AccountID,
			"Type":        t.Type,
			"Amount":      t.Amount,
			"Description": t.Description,
			"IsRecurring": t.IsRecurring,
			"CreatedAt":   t.CreatedAt.Format(time.RFC3339),
		}
	}

	response := map[string]interface{}{
		"transactions": transactionData,
	}

	responseData, _ := json.Marshal(response)

	Send_Message(ws, string(responseData))
}

func GetTransactionById(
	ws *websocket.Conn,
	data json.RawMessage,
	userID string,
) {
	transaction := new(types.Transaction)

	transaction.UserID = userID

	if err := json.Unmarshal(data, &transaction); err != nil {
		Send_Error(ws, InvalidData, err)
		return
	}

	if transaction.UserID != userID {
		Send_Error(ws, "Transaction does not belong to you", nil)
		return
	}

	if err := db.DB.Where("id = ?", transaction.ID).First(&transaction).Error; err != nil {
		Send_Error(ws, "Transaction not found", err)
		return
	}

	transactionData := map[string]interface{}{
		"ID":          transaction.ID,
		"UserID":      transaction.UserID,
		"AccountID":   transaction.AccountID,
		"Type":        transaction.Type,
		"Amount":      transaction.Amount,
		"IsRecurring": transaction.IsRecurring,
	}

	response := map[string]interface{}{
		"Success": transactionData,
	}

	responseData, _ := json.Marshal(response)
	Send_Message(ws, string(responseData))
}

func UpdateTransaction(
	ws *websocket.Conn,
	data json.RawMessage,
	userID string,
) {
	transaction := new(types.Transaction)

	if err := json.Unmarshal(data, &transaction); err != nil {
		Send_Error(ws, InvalidData, err)
	}

	if transaction.UserID != userID {
		Send_Error(ws, "Transaction does not belong to you", nil)
		return
	}

	if err := db.DB.Where("id = ?", transaction.ID).First(&transaction).Error; err != nil {
		Send_Error(ws, "Transaction not found", err)
		return

	}

	if err := db.DB.Save(transaction).Error; err != nil {
		Send_Error(ws, "Failed to save", err)
	}

	transactionData := map[string]interface{}{
		"ID":        transaction.ID,
		"UserID":    transaction.UserID,
		"Type":      transaction.Type,
		"AccountID": transaction.AccountID,
		"Amount":    transaction.Amount,
	}

	response := map[string]interface{}{
		"Success": transactionData,
	}

	responseData, _ := json.Marshal(response)
	Send_Message(ws, string(responseData))
}

func DeleteTransaction(
	ws *websocket.Conn,
	data json.RawMessage,
	userID string,
) {
	transaction := new(types.Transaction)

	transaction.UserID = userID

	if err := json.Unmarshal(data, &transaction); err != nil {
		Send_Error(ws, InvalidData, err)

		return
	}
	if transaction.UserID != userID {
		Send_Error(ws, "Transaction does not belong to you", nil)
		return
	}

	if err := db.DB.Where("id = ?", transaction.ID).First(&transaction).Error; err != nil {
		Send_Error(ws, "Transaction not found", err)
		return
	}

	if err := db.DB.Delete(transaction).Error; err != nil {
		Send_Error(ws, "Failed to delete", err)
		return
	}
	response := map[string]interface{}{
		"Success": "Transaction Deleted",
	}

	responseData, _ := json.Marshal(response)
	Send_Message(ws, string(responseData))
}

func GetAccountIncome(
	ws *websocket.Conn,
	data json.RawMessage,
	userID string,
) {
	transactions := []types.Transaction{}
	account := new(types.Accounts)

	account.ID = requestData.AccountID
	if err := json.Unmarshal(data, &requestData); err != nil {
		Send_Error(ws, InvalidData, err)
		return
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, requestData.AccountID).First(&account).Error; err != nil {
		Send_Error(ws, "Account not found", err)
		return
	}

	start, end := GetDateRange(requestData.DateRange)
	if err := db.DB.Where("account_id = ? and type = ? AND created_at BETWEEN ? AND ?", requestData.AccountID, "Income", start, end).Find(&transactions).Error; err != nil {
		Send_Error(ws, "Could not get transactions", err)
	}

	totalIncome := 0.0
	for _, transaction := range transactions {
		totalIncome += transaction.Amount
	}
	account.Income = totalIncome

	if err := db.DB.Save(account).Error; err != nil {
		Send_Error(ws, "Failed to save", err)
	}
	response := map[string]interface{}{
		"totalIncome": account.Income,
	}

	responseData, _ := json.Marshal(response)
	Send_Message(ws, string(responseData))
}

func GetAccountExpense(
	ws *websocket.Conn,
	data json.RawMessage,
	userID string,
) {
	transactions := []types.Transaction{}
	account := new(types.Accounts)

	account.ID = requestData.AccountID
	if err := json.Unmarshal(data, &requestData); err != nil {
		Send_Error(ws, InvalidData, err)
		return
	}
	if err := db.DB.Where("user_id =? AND id =?", userID, requestData.AccountID).First(&account).Error; err != nil {
		Send_Error(ws, "Account not found", err)
		return
	}

	start, end := GetDateRange(requestData.DateRange)
	if err := db.DB.Where("account_id = ? and type = ? AND created_at BETWEEN ? AND ?", requestData.AccountID, "Expense", start, end).Find(&transactions).Error; err != nil {
		Send_Error(ws, "Could not get transactions", err)
	}
	totalExpense := 0.0
	for _, transaction := range transactions {
		totalExpense += transaction.Amount
	}
	account.Expense = totalExpense

	if err := db.DB.Save(account).Error; err != nil {
		Send_Error(ws, "Failed to Save", err)
	}

	response := map[string]interface{}{
		"totalExpense": account.Expense,
	}

	responseData, _ := json.Marshal(response)
	Send_Message(ws, string(responseData))
}

func GetAccountBalance(
	ws *websocket.Conn,
	data json.RawMessage,
	userID string,
) {
	account := new(types.Accounts)

	account.ID = requestData.AccountID
	if err := json.Unmarshal(data, &requestData); err != nil {
		Send_Error(ws, InvalidData, err)
		return
	}
	if err := db.DB.Where("user_id =? AND id =?", userID, requestData.AccountID).First(&account).Error; err != nil {
		Send_Error(ws, "Account not found", err)
		return
	}

	totalBalance := account.Income - account.Expense
	account.Balance = totalBalance

	if err := db.DB.Save(account).Error; err != nil {
		Send_Error(ws, "Failed to Save", err)
	}

	response := map[string]interface{}{
		"accountBalance": account.Balance,
	}

	responseData, _ := json.Marshal(response)
	Send_Message(ws, string(responseData))
}

func GetAccountsBalance(ws *websocket.Conn, accountID string) error {
	transactions := []types.Transaction{}
	account := new(types.Accounts)

	account.ID = accountID

	if err := db.DB.Where(" id =?", account.ID).First(&account).Error; err != nil {
		Send_Error(ws, "Account not found", err)
		return err
	}

	if err := db.DB.Where("account_id = ?", account.ID).Find(&transactions).Error; err != nil {
		Send_Error(ws, "Could not get transactions", err)
		return err
	}

	totalBalance := account.Income - account.Expense
	account.Balance = totalBalance

	if err := db.DB.Save(account).Error; err != nil {
		Send_Error(ws, "Failed to save", err)
		return err
	}

	return nil
}
