package handlers

import (
	"context"
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
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
		Send_Error(ws, "Invalid form data", err)
		return
	}
	account := &types.Accounts{
		ID:      uuid.NewString(),
		UserID:  userID,
		Type:    req.Type,
		Balance: req.Balance,
		Income:  req.Income,
		Expense: req.Expense,
	}

	var userExists bool
	err := db.DB.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)",
		userID,
	).Scan(&userExists)
	if err != nil || !userExists {
		Send_Error(ws, "User not found", err)
		return
	}

	if _, err := db.DB.Exec(context.Background(),
		`INSERT INTO accounts (
id,
user_id,
type,
balance,
income,
expense
)VALUES ($1, $2, $3, $4, $5, $6)`,
		account.ID,
		account.UserID,
		account.Type,
		account.Balance,
		account.Income,
		account.Expense,
	); err != nil {
		Send_Error(ws, "Failed to Create Account", err)
	}
	response := map[string]any{
		"account": map[string]any{
			"accountID": account.ID,
			"type":      account.Type,
			"balance":   account.Balance,
		},
	}

	responseData, _ := json.Marshal(response)
	Send_Message(ws, string(responseData))
}

func GetAccounts(ws *websocket.Conn, data json.RawMessage, userID string) {
		accounts := []types.Accounts{}
	rows, err := db.DB.Query(
		context.Background(),
		`SELECT id, type, balance, income, expense 
         FROM accounts WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		Send_Error(ws, "Failed to fetch accounts", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var a types.Accounts
		if err := rows.Scan(
			&a.ID, &a.Type, &a.Balance, &a.Income, &a.Expense,
		); err != nil {
			Send_Error(ws, "Data processing error", err)
			return
		}
		accounts = append(accounts, a)
	}

	accountData := make([]map[string]any, 0, len(accounts))
	for _, a := range accounts {
		if err := GetAccountsBalance(ws, a.ID); err != nil {
			Send_Error(ws, "Failed to get balance", err)
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
		Send_Error(ws, "Failed to generate response", err)
		return
	}
	Send_Message(ws, string(responseData))
}

func UpdateAccount(ws *websocket.Conn, data json.RawMessage, userID string) {
	account := new(types.Accounts)
	if err := json.Unmarshal(data, &account); err != nil {
		Send_Error(ws, InvalidData, err)
		return
	}

	if err := db.DB.QueryRow(context.Background(),
		"SELECT * FROM accounts WHERE user_id = $1", userID).Scan(&account); err != nil {
		Send_Error(ws, "Account not found", err)
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
		Send_Error(ws, "Failed to update", err)
	}
	accountData := map[string]any{
		"ID":   account.ID,
		"Type": account.Type,
	}

	response := map[string]any{
		"Success": accountData,
	}

	responseData, _ := json.Marshal(response)

	Send_Message(ws, string(responseData))
}

func DeleteAccount(ws *websocket.Conn, data json.RawMessage, userID string) {
	account := new(types.Accounts)
	if err := json.Unmarshal(data, &account); err != nil {
		Send_Error(ws, InvalidData, err)

		return

	}

	if err := db.DB.QueryRow(context.Background(),
		"SELECT * FROM accounts WHERE user_id = $1", userID).Scan(&account); err != nil {
		Send_Error(ws, "Account not found", err)
		return
	}
	if _, err := db.DB.Exec(context.Background(), `DELETE FROM accounts WHERE id = $1`, account.ID); err != nil {
		Send_Error(ws, "Failed to delete", err)
	}

	accountData := map[string]any{
		"ID":   account.ID,
		"Type": account.Type,
	}

	response := map[string]any{
		"Success": accountData,
	}

	responseData, _ := json.Marshal(response)
	Send_Message(ws, string(responseData))
}
