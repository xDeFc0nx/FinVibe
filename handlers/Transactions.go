package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/xDeFc0nx/NovaoFin/db"
	"github.com/xDeFc0nx/NovaoFin/types"
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
		SendError(ws, MsgInvalidFrequency, nil)
		return nil
	}
	if err := json.Unmarshal(data, &recurring); err != nil {
		SendError(ws, MsgInvalidData, err)
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
		SendError(ws, "Invalid frequency", nil)
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
		SendError(ws, fmt.Sprintf(MsgCreateFailedFmt, "Transaction"), err)
	}
	go func() {
		err := handleRecurringTransaction(ws, recurring, userID, accountID)
		if err != nil {
			SendError(ws, "Failed to handle recurring transactions", err)
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
		SendError(ws, MsgInvalidData, err)
		return
	}

	if transaction.AccountID == "" {
		SendError(ws, MsgMissingAccountID, nil)
		return
	}
	if _, err := db.DB.Exec(context.Background(), `
SELECT EXISTS (SELECT 1 FROM accounts WHERE id = $1 AND user_id = $2)
		`, transaction.AccountID, userID); err != nil {
		SendError(ws, MsgAccountNotFound, err)
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
		SendError(ws, fmt.Sprintf(MsgCreateFailedFmt, "Transaction"), err)
	}
	if transaction.IsRecurring {

		var inputData map[string]any
		if err := json.Unmarshal(data, &inputData); err != nil {

			SendError(ws, "Failed to parse recurring frequency", err)
			return
		}

		if freq, ok := inputData["Frequency"].(string); ok && freq != "" {
			recurring.Frequency = freq
		} else {

			SendError(ws, fmt.Sprintf(MsgMissingFieldFmt, "Frequency"), nil)
			return
		}
		if err := CreateRecurring(ws, data, transaction.ID, account.ID, userID, transaction.Amount); err != nil {
			return
		}
	}
	if err := GetAccountBalance(ws, transaction.AccountID); err != nil {
		SendError(ws, fmt.Sprintf(MsgFetchFailedFmt, "Account Balance"), err)
	}
	if err := db.DB.QueryRow(context.Background(), `
		SELECT income, expense, balance,
		FROM accounts
		WHERE id = $1
		`, transaction.AccountID).Scan(&account.Income, &account.Expense, &account.Balance); err != nil {
		SendError(ws, MsgAccountNotFound, err)
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
		"AccountData": map[string]any{
			"Income":   account.Income,
			"Expenses": account.Expense,
			"Balance":  account.Balance,
		},
	}

	responseData, _ := json.Marshal(response)

	SendMessage(ws, string(responseData))
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
			SendError(ws, fmt.Sprintf(MsgCreateFailedFmt, "Transactions"), err)
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
			SendError(ws, fmt.Sprintf(MsgUpdateFailedFmt, "Recurring"), err)
		}
	}
}

func GetTransactions(ws *websocket.Conn, data json.RawMessage, userID string) {

	if err := json.Unmarshal(data, &requestData); err != nil {
		SendError(ws, MsgInvalidData, err)
		return
	}

	if requestData.DateRange == "" {
		SendError(ws, requestData.DateRange, nil)
	}
	if _, err := db.DB.Exec(context.Background(), `
SELECT EXISTS (SELECT 1 FROM accounts WHERE id = $1 AND user_id = $2)
		`, requestData.AccountID, userID); err != nil {
		SendError(ws, MsgAccountNotFound, err)
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
		SendError(ws, fmt.Sprintf(MsgFetchFailedFmt, "Transactions"), err)
	}

	defer rows.Close()
	transactions, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.Transaction])
	if err != nil {
		slog.Error("failed", slog.String("err", err.Error()))
		SendError(ws, MsgCollectRowsFailed, err)
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

	SendMessage(ws, string(responseData))
}

func GetTransactionById(
	ws *websocket.Conn,
	data json.RawMessage,
	userID string,
) {
	transaction := new(types.Transaction)

	transaction.UserID = userID

	if err := json.Unmarshal(data, &transaction); err != nil {
		SendError(ws, MsgInvalidData, err)
		return
	}

	if transaction.UserID != userID {
		SendError(ws, MsgUnauthorized, nil)
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
		SendError(ws, fmt.Sprintf(MsgFetchFailedFmt, "Transactions"), err)
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
	SendMessage(ws, string(responseData))
}

func UpdateTransaction(
	ws *websocket.Conn,
	data json.RawMessage,
	userID string,
) {
	transaction := new(types.Transaction)
	account := new(types.Accounts)
	if err := json.Unmarshal(data, &transaction); err != nil {
		SendError(ws, MsgInvalidData, err)
	}

	if transaction.UserID != userID {
		SendError(ws, MsgUnauthorized, nil)
		return
	}
	if _, err := db.DB.Query(context.Background(), `
SELECT EXISTS transactions WHERE id = $1 AND user_id = $2
		`, transaction.ID, userID); err != nil {
		SendError(ws, MsgTransactionNotFound, err)

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
		SendError(ws, fmt.Sprintf(MsgFetchFailedFmt, "Transactions"), err)
	}

	if err := GetAccountBalance(ws, transaction.AccountID); err != nil {
		SendError(ws, fmt.Sprintf(MsgFetchFailedFmt, "Account Balance"), err)
	}
	if err := db.DB.QueryRow(context.Background(), `
		SELECT income, expense, balance,
		FROM accounts
		WHERE id = $1
		`, transaction.AccountID).Scan(&account.Income, &account.Expense, &account.Balance); err != nil {
		SendError(ws, MsgAccountNotFound, err)
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
		"AccountData": map[string]any{
			"Income":  account.Income,
			"Expense": account.Expense,
			"Balance": account.Balance,
		},
	}
	responseData, _ := json.Marshal(response)
	SendMessage(ws, string(responseData))

}

func DeleteTransaction(ws *websocket.Conn, data json.RawMessage,
	userID string,
) {
	transaction := new(types.Transaction)
	account := new(types.Accounts)
	transaction.UserID = userID

	if err := json.Unmarshal(data, &transaction); err != nil {
		SendError(ws, MsgInvalidData, err)

		return
	}
	if transaction.UserID != userID {
		SendError(ws, MsgUnauthorized, nil)
		return
	}
	if _, err := db.DB.Exec(context.Background(), `
SELECT EXISTS transactions WHERE id = $1 AND user_id = $2
		`, transaction.ID, userID); err != nil {
		SendError(ws, MsgTransactionNotFound, err)
	}
	if _, err := db.DB.Exec(context.Background(), `
			DELETE transactions WHERE id = $1 AND user_id = $2
			`, transaction.ID, userID); err != nil {
		SendError(ws, fmt.Sprintf(MsgDeleteFailedFmt, "Transaction"), err)
	}
	if err := GetAccountBalance(ws, transaction.AccountID); err != nil {
		SendError(ws, fmt.Sprintf(MsgFetchFailedFmt, "Account Balance"), err)
	}
	if err := db.DB.QueryRow(context.Background(), `
		SELECT income, expense, balance,
		FROM accounts
		WHERE id = $1
		`, transaction.AccountID).Scan(&account.Income, &account.Expense, &account.Balance); err != nil {
		SendError(ws, MsgAccountNotFound, err)
	}
	response := map[string]any{
		"Success": "Transaction Deleted",
		"AccountData": map[string]any{
			"Income":  account.Income,
			"Expense": account.Expense,
			"Balance": account.Balance,
		},
	}

	responseData, _ := json.Marshal(response)
	SendMessage(ws, string(responseData))
}
