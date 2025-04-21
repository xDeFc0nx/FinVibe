package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"log/slog"
	"sync"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func CreateBudget(ws *websocket.Conn, data json.RawMessage, userID string) {
	budget := new(types.Budget)

	budget.ID = uuid.NewString()
	budget.UserID = userID

	type requestData struct {
		AccountID   string  `json:"accountId"`
		TotalSpent  float64 `json:"totalSpent"`
		LimitAmount float64 `json:"limit"`
		Description string  `json:"description"`
	}

	var req requestData

	if err := json.Unmarshal(data, &req); err != nil {
		SendError(ws, MsgInvalidData, err)
	}

	if budget.ID == "" {
		SendError(ws, fmt.Sprintf(MsgMissingFieldFmt, ("Bugdet ID")), nil)
	}

	if _, err := db.DB.Exec(context.Background(), `	
	SELECT FROM accounts WHERE id = $1 AND user_id = $2
		`, req.AccountID, userID); err != nil {
		SendError(ws, MsgAccountNotFound, err)
	}

	if _, err := db.DB.Exec(context.Background(), `
		INSERT INTO budgets (id, user_id, account_id, limit, total_spent, description)
		VALUES ($1, $2, $3, $4, $5, $6)

		`,
		budget.ID,
		budget.UserID,
		req.AccountID,
		req.LimitAmount,
		req.TotalSpent,
		req.Description,
	); err != nil {
		SendError(ws, MsgWebSocketCreationFailed, err)
	}
	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID":          budget.ID,
			"Amount":      budget.LimitAmount,
			"Description": budget.Description,
		},
	}

	responseData, _ := json.Marshal(response)

	SendMessage(ws, string(responseData))
}

func GetBudgets(ws *websocket.Conn, data json.RawMessage, userID string) {
	budgets := []types.Budget{}

	type requestData struct {
		AccountID string `json:"AccountID"`
	}
	var req requestData
	if _, err := db.DB.Exec(context.Background(), `	
	SELECT FROM budgets WHERE id = $1 AND user_id = $2
		`, req.AccountID, userID); err != nil {
		SendError(ws, MsgAccountNotFound, err)
	}

	var wg sync.WaitGroup
	for i := range budgets {
		wg.Add(1)
		go func(a *types.Budget) {
			defer wg.Done()
			if err := GetBudgetCal(ws, a.ID); err != nil {
				slog.Error(
					"Failed to get budget",
					slog.String("error", err.Error()),
				)
				SendError(ws, fmt.Sprintf(MsgFetchFailedFmt, "Bugdet"), err)
			}
		}(&budgets[i])
	}
	wg.Wait()
	if _, err := db.DB.Exec(context.Background(), `	
	SELECT FROM budgets WHERE user_id = $1
		`, userID); err != nil {
		SendError(ws, MsgBudgetNotFound, err)
	}
	budgetsData := make([]map[string]interface{}, len(budgets))

	for i, a := range budgets {
		budgetsData[i] = map[string]interface{}{
			"ID":            a.ID,
			"UserID":        a.UserID,
			"AccountID":     a.ID,
			"Budgets Limit": float64(a.LimitAmount),
			"Budgets Spend": float64(a.TotalSpent),
		}
	}

	response := map[string]interface{}{
		"Success": "Fetched budgets",
		"budgets": budgetsData,
	}

	responseData, _ := json.Marshal(response)

	SendMessage(ws, string(responseData))
}

func UpdateBudget(ws *websocket.Conn, data json.RawMessage, userID string) {
	type requestData struct {
		AccountID   string  `json:"accountId"`
		TotalSpent  float64 `json:"totalSpent"`
		LimitAmount float64 `json:"limit"`
		Description string  `json:"description"`
	}

	var req requestData
	budget := new(types.Budget)

	if err := json.Unmarshal(data, req); err != nil {
		SendError(ws, MsgInvalidData, err)
	}
	if budget.ID == "" {
		SendError(ws, "ID is required", nil)
	}

	if _, err := db.DB.Exec(context.Background(), `	
	SELECT FROM budgets WHERE id = $1 AND user_id = $2
		`, budget.ID, userID); err != nil {
		SendError(ws, MsgBudgetNotFound, err)
	}
	if _, err := db.DB.Exec(context.Background(), `
	INSERT INTO budgets(id, user_id, account_id, limit, total_spent, description)
		VALUES($1, $2, $3, $4, $5, $6)
		`,
		budget.ID,
		budget.UserID,
		req.AccountID,
		req.LimitAmount,
		req.TotalSpent,
		req.Description,
	); err != nil {
		SendError(ws, fmt.Sprintf(MsgCreateFailedFmt, "Bugdet"), err)
	}

	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID":     budget.ID,
			"UserID": budget.UserID,
			"Amount": budget.LimitAmount,
		},
	}

	responseData, _ := json.Marshal(response)

	SendMessage(ws, string(responseData))

}
func DeleteBudget(ws *websocket.Conn, data json.RawMessage, userID string) {
	budget := new(types.Budget)

	if err := json.Unmarshal(data, budget); err != nil {
		SendError(ws, MsgInvalidData, err)
	}

	if budget.ID == "" {
		SendError(ws, "ID is required", nil)
	}

	if _, err := db.DB.Exec(context.Background(), `
    DELETE FROM budgets WHERE id = $1 AND user_id = $2
		`,
		budget.ID,
		userID,
	); err != nil {
		SendError(ws, fmt.Sprintf(MsgDeleteFailedFmt, "Budget"), err)
	}

	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID": budget.ID,
		},
	}

	responseData, _ := json.Marshal(response)

	SendMessage(ws, string(responseData))
}

func GetBudgetCal(ws *websocket.Conn, accountID string) error {
	transactions := []types.Transaction{}
	account := new(types.Accounts)
	budget := new(types.Budget)

	account.ID = accountID

	if _, err := db.DB.Exec(context.Background(), `	
	SELECT FROM accounts WHERE id = $1
		`, accountID); err != nil {
		SendError(ws, MsgBudgetNotFound, err)
	}
	if _, err := db.DB.Exec(context.Background(), `
 SELECT FROM transactions WHERE account_id = $1
		`, accountID); err != nil {
		SendError(ws, MsgTransactionNotFound, err)
	}
	totalBalance := float64(0)
	for _, t := range transactions {
		if t.AccountID == accountID {
			totalBalance += t.Amount
		}
	}

	budget.TotalSpent = totalBalance
	if _, err := db.DB.Exec(context.Background(), `
INSERT INTO budgets 
		total_spent = $1
		WHERE account_id = $2
		`, budget.TotalSpent, accountID); err != nil {
		SendError(ws, fmt.Sprintf(MsgUpdateFailedFmt, "Budget"), err)
	}
	return nil
}
