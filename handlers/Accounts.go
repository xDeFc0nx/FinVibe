package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/xDeFc0nx/NovaoFin/db"
	"github.com/xDeFc0nx/NovaoFin/types"
)

func CreateAccount(ws *websocket.Conn, data json.RawMessage, userID string) {
	type Request struct {
		Type    string  `json:"type"`
		Balance float64 `json:"balance"`
		Income  float64 `json:"income"`
		Expense float64 `json:"expense"`
	}
	var req Request



	if err := json.Unmarshal(data, &req); err != nil {
		SendError(ws, MsgInvalidData, err)
		return
	}
	account := &types.Accounts{}
	account.ID = uuid.NewString()
	account.Type = req.Type
	account.Balance = req.Balance
	account.Income = req.Income
	account.Expense = req.Expense
	var userExists bool
	err := db.DB.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)",
		userID,
	).Scan(&userExists)
	if err != nil || !userExists {
		SendError(ws, MsgUserNotFound, err)
		return
	}
	if _, err := db.DB.Exec(context.Background(),
		`INSERT INTO accounts (
        id,
        user_id,
        type,
        balance,
        income,
        expense,
        created_at,
        updated_at
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		account.ID,
		userID,
		account.Type,
		account.Balance,
		account.Income,
		account.Expense,
		time.Now().UTC(),
		time.Now().UTC(),
	); err != nil {
		SendError(ws, MsgCreateFailedFmt, err)
	}
	response := map[string]any{
		"account": map[string]any{
			"accountID": account.ID,
			"type":      account.Type,
			"balance":   account.Balance,
		},
	}

	responseData, _ := json.Marshal(response)
	SendMessage(ws, string(responseData))
}

func GetAccounts(ws *websocket.Conn, data json.RawMessage, userID string) {
	rows, err := db.DB.Query(
		context.Background(),
		`SELECT income, expense, balance, id, user_id, type, created_at, updated_at
     FROM accounts
		 WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		SendError(ws, fmt.Sprintf(MsgFetchFailedFmt, "Account"), err)
		return
	}

	defer rows.Close()
	accounts, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.Accounts])
	if err != nil {
		SendError(ws, MsgCollectRowsFailed, err)
	}
	accountData := make([]map[string]any, 0, len(accounts))
	for _, a := range accounts {
		if err := GetAccountBalance(ws, a.ID); err != nil {
			SendError(ws, MsgAccountNotFound, err)

			continue
		}
		accountData = append(accountData, map[string]any{
			"id":      a.ID,
			"type":    a.Type,
			"balance": a.Balance,
			"income":  a.Income,
			"expense": a.Expense,
		})
	}

	response := map[string]any{"accounts": accountData}
	responseData, err := json.Marshal(response)
	if err != nil {
		SendError(ws, MsgGenerateResponseFailed, err)
		return
	}
	SendMessage(ws, string(responseData))
}

func GetAccount(ws *websocket.Conn, data json.RawMessage, userID string) {
	type Request struct {
		AccountID string `json:"AccountID"`
	}
	var req Request

	if err := json.Unmarshal(data, &req); err != nil {
		SendError(ws, MsgInvalidData, err)
		return
	}
	account := new(types.Accounts)
	if err := db.DB.QueryRow(context.Background(), `
  SELECT 1 id, type, income, expense, balance
	FROM accounts
	WHERE id = $1
				`, req.AccountID).Scan(&account.ID, &account.Type, &account.Income, &account.Expense, &account.Balance); err != nil {
		SendError(ws, MsgAccountNotFound, err)
	}
	accountData := map[string]any{
		"id":      account.ID,
		"type":    account.Type,
		"balance": account.Balance,
		"income":  account.Income,
		"expense": account.Expense,
	}

	responseData, _ := json.Marshal(accountData)

	SendMessage(ws, string(responseData))

}
func UpdateAccount(ws *websocket.Conn, data json.RawMessage, userID string) {
	account := new(types.Accounts)
	if err := json.Unmarshal(data, &account); err != nil {
		SendError(ws, MsgInvalidData, err)
		return
	}

	if err := db.DB.QueryRow(context.Background(),
		"SELECT * FROM accounts WHERE user_id = $1", userID).Scan(&account); err != nil {
		SendError(ws, MsgAccountNotFound, err)
		return
	}
	if _, err := db.DB.Exec(context.Background(), `
		UPDATE accounts SET
		type = $1,
		balance = $2,
		income = $3,
		expense = $4
		WHERE id = $5 AND user_id = $6`,
		account.Type,
		account.Balance,
		account.Income,
		account.Expense,
		account.ID,
		account.UserID,
	); err != nil {
		SendError(ws, fmt.Sprintf(MsgUpdateFailedFmt, "Account"), err)
	}
	accountData := map[string]any{
		"ID":   account.ID,
		"Type": account.Type,
	}

	response := map[string]any{
		"Success": accountData,
	}

	responseData, _ := json.Marshal(response)

	SendMessage(ws, string(responseData))
}

func DeleteAccount(ws *websocket.Conn, data json.RawMessage, userID string) {
	account := new(types.Accounts)
	if err := json.Unmarshal(data, &account); err != nil {
		SendError(ws, MsgInvalidData, err)

		return

	}

	if err := db.DB.QueryRow(context.Background(),
		"SELECT * FROM accounts WHERE user_id = $1", userID).Scan(&account); err != nil {
		SendError(ws, fmt.Sprintf(MsgFetchFailedFmt, "Account"), err)
		return
	}
	if _, err := db.DB.Exec(context.Background(), `DELETE FROM accounts WHERE id = $1`, account.ID); err != nil {
		SendError(ws, fmt.Sprintf(MsgDeleteFailedFmt, "Account"), err)
	}

	accountData := map[string]any{
		"ID":   account.ID,
		"Type": account.Type,
	}

	response := map[string]any{
		"Success": accountData,
	}

	responseData, _ := json.Marshal(response)
	SendMessage(ws, string(responseData))
}
func GetAccountBalance(
	ws *websocket.Conn,
	AccountID string,
) error {
	account := new(types.Accounts)

	account.ID = AccountID

	if err := db.DB.QueryRow(context.Background(), `
SELECT 1 id, type, income, expense, balance
	FROM accounts
		WHERE id = $1 
				`, AccountID).Scan(&account.ID, &account.Type, &account.Income, &account.Expense, &account.Balance); err != nil {
		SendError(ws, MsgAccountNotFound, err)
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
		AccountID,
		start,
		end,
		incType,
	)
	if err != nil {
		SendError(ws, fmt.Sprintf(MsgFetchFailedFmt, "Transactions"), err)
	}

	defer rows.Close()
	transactions, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.Transaction])
	if err != nil {
		SendError(ws, MsgCollectRowsFailed, err)
	}

	totalIncome := 0.0
	for _, transaction := range transactions {
		totalIncome += transaction.Amount
	}

	if err := db.DB.QueryRow(context.Background(), `
SELECT 1 id, type, income, expense, balance
	FROM accounts
		WHERE id = $1
				`, AccountID).Scan(&account.ID, &account.Type, &account.Income, &account.Expense, &account.Balance); err != nil {
		SendError(ws, MsgAccountNotFound, err)
	}

	account.Income = totalIncome
	expType := "Expense"
	rows, err = db.DB.Query(context.Background(), `
SELECT amount, id, user_id, account_id, type, description, is_recurring, created_at, updated_at
		FROM transactions
		WHERE account_id = $1 
		AND created_at BETWEEN $2 AND $3
		AND type = $4
		ORDER BY created_at DESC`,
		AccountID,
		start,
		end,
		expType,
	)
	if err != nil {
		SendError(ws, fmt.Sprintf(MsgFetchFailedFmt, "Transactions"), err)
	}

	defer rows.Close()
	transactions, err = pgx.CollectRows(rows, pgx.RowToStructByName[types.Transaction])
	if err != nil {
		SendError(ws, MsgCollectRowsFailed, err)
	}

	totalExpense := 0.0
	for _, transaction := range transactions {
		totalExpense += transaction.Amount
	}

	totalBalance := account.Income - account.Expense
	if _, err := db.DB.Exec(context.Background(), `
		UPDATE accounts SET
		income = $1,
		expense = $2,
		balance = $3
		WHERE id = $4
	
		`, totalIncome, totalExpense, totalBalance, AccountID); err != nil {
		SendError(ws, fmt.Sprintf(MsgUpdateFailedFmt, "Account Balance"), err)
	}
	return nil
}
