package handlers

import (
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func CreateBudget(ws *websocket.Conn, data json.RawMessage, userID string) {
	budget := new(types.Budget)

	budget.ID = uuid.NewString()
	budget.UserID = userID

	var requestData struct {
		AccountID string `json:"AccountID"`
	}

	if err := json.Unmarshal(data, budget); err != nil {
		Message(ws, InvalidData, err)
	}

	if budget.ID == "" {
		Message(ws, "ID is required", nil)
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, requestData.AccountID).First(&budget.Account).Error; err != nil {
		Message(ws, "Account not found", err)
	}

	if err := db.DB.Create(&budget).Error; err != nil {
		Message(ws, "Failed to create budget", err)
	}
	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID":          budget.ID,
			"Amount":      budget.Limit,
			"Description": budget.Description,
		},
	}

	responseData, _ := json.Marshal(response)

	if err := ws.WriteMessage(websocket.TextMessage, responseData); err != nil {
		slog.Error("Failed to send message", slog.String("error", err.Error()))
	}
}

func GetBudgets(ws *websocket.Conn, data json.RawMessage, userID string) {
	budgets := []types.Budget{}

	var requestData struct {
		AccountID string `json:"AccountID"`
	}

	if err := db.DB.Where("user_id =? AND account_id = ?", userID, requestData.AccountID).Find(&budgets).Error; err != nil {
		Message(ws, "Budgets not found", err)
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
			}
		}(&budgets[i])
	}
	wg.Wait()
	if err := db.DB.Where("user_id = ?", userID).Find(&budgets).Error; err != nil {
		Message(ws, "Budgets not found", err)
	}

	budgetsData := make([]map[string]interface{}, len(budgets))

	for i, a := range budgets {
		budgetsData[i] = map[string]interface{}{
			"ID":            a.ID,
			"UserID":        a.UserID,
			"AccountID":     a.ID,
			"Budgets Limit": float64(a.Limit),
			"Budgets Spend": float64(a.TotalSpent),
		}
	}

	// Package the response
	response := map[string]interface{}{
		"Success": "Fetched budgets",
		"budgets": budgetsData,
	}

	responseData, _ := json.Marshal(response)
	if err := ws.WriteMessage(websocket.TextMessage, responseData); err != nil {
		slog.Error("Failed to send message", slog.String("error", err.Error()))
	}
}

func UpdateBudget(ws *websocket.Conn, data json.RawMessage, userID string) {
	budget := new(types.Budget)

	if err := json.Unmarshal(data, budget); err != nil {
		Message(ws, InvalidData, err)
	}
	if budget.ID == "" {
		Message(ws, "ID is required", nil)
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, budget.ID).First(&budget).Error; err != nil {
		Message(ws, "Budget not found", err)
	}

	if err := db.DB.Save(budget).Error; err != nil {
		Message(ws, "Failed to save", err)
	}

	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID":     budget.ID,
			"UserID": budget.UserID,
			"Amount": budget.Limit,
		},
	}

	responseData, _ := json.Marshal(response)

	if err := ws.WriteMessage(websocket.TextMessage, responseData); err != nil {
		slog.Error("Failed to send message", slog.String("error", err.Error()))
	}
}

func DeleteBudget(ws *websocket.Conn, data json.RawMessage, userID string) {
	budget := new(types.Budget)

	if err := json.Unmarshal(data, budget); err != nil {
		Message(ws, InvalidData, err)
	}

	if budget.ID == "" {
		Message(ws, "ID is required", nil)
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, budget.ID).Delete(&budget).Error; err != nil {
		Message(ws, "Failed to delete budget", err)
	}

	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID": budget.ID,
		},
	}

	responseData, _ := json.Marshal(response)

	if err := ws.WriteMessage(websocket.TextMessage, responseData); err != nil {
		slog.Error("Failed to send message", slog.String("error", err.Error()))
	}
}

func GetBudgetCal(ws *websocket.Conn, accountID string) error {
	transactions := []types.Transaction{}
	account := new(types.Accounts)
	budget := new(types.Budget)

	account.ID = accountID

	if err := db.DB.Where(" id =?", account.ID).First(&account).Error; err != nil {
		Message(ws, "Account not found", err)
		return err
	}

	if err := db.DB.Where("account_id = ?", account.ID).Find(&transactions).Error; err != nil {

		Message(ws, "Could not fetch transactions", err)
		return err
	}

	totalBalance := float64(0)
	for _, t := range transactions {
		if t.AccountID == accountID {
			totalBalance += t.Amount
		}
	}

	budget.TotalSpent = totalBalance

	if err := db.DB.Save(budget).Error; err != nil {
		Message(ws, "Failed to save", err)
		return err
	}

	return nil
}
