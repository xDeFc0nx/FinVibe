package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
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
	if _, err := db.DB.Exec(context.Background(), `

INSERT INTO recurrings (id, transaction_id, amount, frequency, start_date, next_date)
			VALUES($1, $2, $3, $4, $5, $6)
			`,
		recurring.ID,
		recurring.TransactionID,
		recurring.Amount,
		recurring.Frequency,
		recurring.StartDate,
		recurring.NextDate,
	); err != nil {
		Send_Error(ws, "Failed to create recurring transaction", err)
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
	if _, err := db.DB.Exec(context.Background(), `
SELECT EXISTS (SELECT 1 FROM accounts WHERE id = $1 AND user_id = $2)
		`, transaction.AccountID, userID); err != nil {
		Send_Error(ws, "Account not found", err)
	}
	if _, err := db.DB.Exec(context.Background(), `
 INSERT INTO transactions (id, user_id, account_id, type, description, amount, is_recurring, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
	
		`,
		transaction.ID,
		transaction.UserID,
		transaction.AccountID,
		transaction.Type,
		transaction.Description,
		transaction.Amount,
		transaction.IsRecurring,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	); err != nil {
		Send_Error(ws, "Failed to create transaction", err)
	}
	if transaction.IsRecurring {

		var inputData map[string]any
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
		if err := CreateRecurring(ws, data, transaction.ID, account.ID, userID, transaction.Amount); err != nil {
			return
		}
	}

	response := map[string]any{
		"transaction": map[string]any{
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
		if _, err := db.DB.Exec(context.Background(), `
INSERT INTO transactions (id, user_id, account_id, amount, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6)
			`,
			newTransaction.ID,
			newTransaction.UserID,
			newTransaction.AccountID,
			newTransaction.Amount,
			newTransaction.CreatedAt,
			newTransaction.UpdatedAt,
		); err != nil {
			Send_Error(ws, "Failed to create new recurring transaction", err)
		}
		switch recurring.Frequency {
		case "Daily":
			recurring.NextDate = recurring.NextDate.Add(24 * time.Hour)
		case "Weekly":
			recurring.NextDate = recurring.NextDate.Add(7 * 24 * time.Hour)
		case "Monthly":
			recurring.NextDate = recurring.NextDate.AddDate(0, 1, 0)
		}
		if _, err := db.DB.Exec(context.Background(), `
UPDATE recurrings
			SET recurring.next_date = $1
			WHERE transaction_id = $2
			`, recurring.NextDate, recurring.TransactionID); err != nil {
			Send_Error(ws, "Failed to update recurring next date", err)
		}
	}
}

func GetTransactions(ws *websocket.Conn, data json.RawMessage, userID string) {

	if err := json.Unmarshal(data, &requestData); err != nil {
		Send_Error(ws, InvalidData, err)
		return
	}

	if requestData.DateRange == "" {
		Send_Error(ws, "Date Range is Required", nil)
	}
	if _, err := db.DB.Exec(context.Background(), `
SELECT EXISTS (SELECT 1 FROM accounts WHERE id = $1 AND user_id = $2)
		`, requestData.AccountID, userID); err != nil {
		Send_Error(ws, "Account not found", err)
	}

	start, end := GetDateRange(requestData.DateRange)

	rows, err := db.DB.Query(context.Background(), `
SELECT amount, id, user_id, account_id, type, description, is_recurring, created_at, updated_at
		FROM transactions
		WHERE account_id = $1 AND created_at BETWEEN $2 AND $3
		ORDER BY created_at DESC`,
		requestData.AccountID,
		start,
		end,
	)
	if err != nil {
		Send_Error(ws, "failed to get trnasaction", err)
	}

	defer rows.Close()
	transactions := []types.Transaction{}
	transactions, err = pgx.CollectRows(rows, pgx.RowToStructByName[types.Transaction])
	if err != nil {
		slog.Error("failed", slog.String("err", err.Error()))
		Send_Error(ws, "failed to collect rows", err)
	}
	transactionData := make([]map[string]any, len(transactions))

	for i, t := range transactions {
		transactionData[i] = map[string]any{
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

	response := map[string]any{
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

	err := db.DB.QueryRow(context.Background(), `
		SELECT *
		FROM transactions
		WHERE id = $1 
	  `,
		transaction.ID).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.AccountID,
		&transaction.Type,
		&transaction.Description,
		&transaction.IsRecurring,
		&transaction.CreatedAt,
	)
	if err != nil {
		Send_Error(ws, "failed to get transactions", err)
	}

	transactionData := map[string]any{
		"ID":          transaction.ID,
		"UserID":      transaction.UserID,
		"AccountID":   transaction.AccountID,
		"Type":        transaction.Type,
		"Amount":      transaction.Amount,
		"IsRecurring": transaction.IsRecurring,
	}

	response := map[string]any{
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
	if _, err := db.DB.Query(context.Background(), `
SELECT EXISTS transactions WHERE id = $1 AND user_id = $2
		`, transaction.ID, userID); err != nil {
		Send_Error(ws, "Transaction not found", err)

	}
	if _, err := db.DB.Exec(context.Background(), `
UPDATE transactions  SET 
		type = $1
		amount = $2,
		description = $3,
		is_recurring = $4,
		updated_at $5,
		WHERE user_id = $6 AND id = $7
		`,
		transaction.Type,
		transaction.Amount,
		transaction.Description,
		transaction.IsRecurring,
		time.Now().UTC(),
		userID,
		transaction.ID,
	); err != nil {
		Send_Error(ws, "Failed to update trnasaction", err)
	}
	transactionData := map[string]any{
		"ID":        transaction.ID,
		"UserID":    transaction.UserID,
		"Type":      transaction.Type,
		"AccountID": transaction.AccountID,
		"Amount":    transaction.Amount,
	}
	response := map[string]any{
		"Success": transactionData,
	}
	responseData, _ := json.Marshal(response)
	Send_Message(ws, string(responseData))

}

func DeleteTransaction(ws *websocket.Conn, data json.RawMessage,
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
	if _, err := db.DB.Exec(context.Background(), `
SELECT EXISTS transactions WHERE id = $1 AND user_id = $2
		`, transaction.ID, userID); err != nil {
		Send_Error(ws, "Trnasaction not found", err)
	}
	if _, err := db.DB.Exec(context.Background(), `
			DELETE transactions WHERE id = $1 AND user_id = $2
			`, transaction.ID, userID); err != nil {
		Send_Error(ws, "Failed to delete trnasaction", err)
	}
	response := map[string]any{
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

	if err := db.DB.QueryRow(context.Background(), `
SELECT 1 id, type, income, expense, balance
	FROM accounts
		WHERE id = $1 AND user_id = $2
				`, account.ID, userID).Scan(&account.ID, &account.Type, &account.Income, &account.Expense, &account.Balance); err != nil {
		Send_Error(ws, "Account not found", err)
	}
  incType := "Income"
	start, end := GetDateRange(requestData.DateRange)
	rows, err := db.DB.Query(context.Background(), `
SELECT amount, id, user_id, account_id, type, description, is_recurring, created_at, updated_at
		FROM transactions
		WHERE account_id = $1 
		AND created_at BETWEEN $2 AND $3
		AND type = $4
		ORDER BY created_at DESC`,
		requestData.AccountID,
		start,
		end,
		incType,
	)
	if err != nil {
		Send_Error(ws, "failed to get transactions", err)
	}

	defer rows.Close()
	transactions, err = pgx.CollectRows(rows, pgx.RowToStructByName[types.Transaction])
	if err != nil {
		Send_Error(ws, "failed to collect rows", err)
	}

	totalIncome := 0.0
	for _, transaction := range transactions {
		totalIncome += transaction.Amount
	}
	account.Income = totalIncome

	if _, err := db.DB.Exec(context.Background(), `
		UPDATE accounts SET
		income = $1
		WHERE id = $2 AND user_id = $3
	
		`, totalIncome, requestData.AccountID, userID); err != nil {
		Send_Error(ws, "Failed to update account income", err)
	}
	response := map[string]any{
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
	if err := db.DB.QueryRow(context.Background(), `
SELECT 1 id, type, income, expense, balance
	FROM accounts
		WHERE id = $1 AND user_id = $2
				`, account.ID, userID).Scan(&account.ID, &account.Type, &account.Income, &account.Expense, &account.Balance); err != nil {
		Send_Error(ws, "Account not found", err)
	}
	expType := "Expense"
	start, end := GetDateRange(requestData.DateRange)
	rows, err := db.DB.Query(context.Background(), `
SELECT amount, id, user_id, account_id, type, description, is_recurring, created_at, updated_at
		FROM transactions
		WHERE account_id = $1 
		AND created_at BETWEEN $2 AND $3
		AND type = $4
		ORDER BY created_at DESC`,
		requestData.AccountID,
		start,
		end,
		expType,
	)
	if err != nil {
		slog.Error("Failed to get transactions", slog.String("error", err.Error()))
		Send_Error(ws, "failed to get transactions", err)
	}

	defer rows.Close()
	transactions, err = pgx.CollectRows(rows, pgx.RowToStructByName[types.Transaction])
	if err != nil {
		Send_Error(ws, "failed to collect rows", err)
	}

totalExpense := 0.0
	for _, transaction := range transactions {
		totalExpense += transaction.Amount
	}
	if _, err := db.DB.Exec(context.Background(), `
		UPDATE accounts SET
		expense = $1
		where id = $2 AND user_id = $3
	
		`, totalExpense, requestData.AccountID, userID); err != nil {
		Send_Error(ws, "Failed to update account Expense", err)
	}

	response := map[string]any{
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

	if err := json.Unmarshal(data, &requestData); err != nil {
		Send_Error(ws, InvalidData, err)
		return
	}

	if err := db.DB.QueryRow(context.Background(), `
SELECT *
		FROM accounts
		WHERE id = $1 AND user_id = $2
				`, requestData.AccountID, userID).Scan(
		&account.Income,
		&account.Expense,
		&account.Balance,
		&account.ID,
		&account.UserID,
		&account.Type,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		Send_Error(ws, "Account not found", err)
	}
	totalBalance := account.Income - account.Expense
	if _, err := db.DB.Exec(context.Background(), `
		UPDATE accounts SET
		balance = $1
		WHERE id = $2 AND user_id = $3
	
		`, totalBalance, requestData.AccountID, userID); err != nil {
		Send_Error(ws, "Failed to update account income", err)
	}
	response := map[string]any{
		"accountBalance": account.Balance,
	}

	responseData, _ := json.Marshal(response)
	Send_Message(ws, string(responseData))
}

func GetAccountsBalance(ws *websocket.Conn, AccountID string) error {
	account := new(types.Accounts)

	if err := db.DB.QueryRow(context.Background(), `
	SELECT income, expense, balance, id, user_id, type, created_at, updated_at
	FROM accounts
		WHERE id = $1
				`, AccountID).Scan(
		&account.Income,
		&account.Expense,
		&account.Balance,
		&account.ID,
		&account.UserID,
		&account.Type,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		Send_Error(ws, "Account not found", err)
	}
	totalBalance := account.Income - account.Expense
	if _, err := db.DB.Exec(context.Background(), `
		UPDATE accounts SET
		balance = $1
		WHERE id = $2
	
		`, totalBalance, AccountID); err != nil {
		Send_Error(ws, "Failed to update account balance", err)
	}

	return nil
}
