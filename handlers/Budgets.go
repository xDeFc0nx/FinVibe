package handlers

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/xDeFc0nx/logger-go-pkg"

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
		Message(ws, InvalidData)
	}

	if budget.ID == "" {
		Message(ws, "Erro: ID is required")
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, requestData.AccountID).First(&budget.Account).Error; err != nil {
		Message(ws, "Error: Account not found")
	}

	if err := db.DB.Create(&budget).Error; err != nil {
		Message(ws, "Error: Failed to create budget")
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
		logger.Error("%s", err.Error())
	}
}

func GetBudgets(ws *websocket.Conn, data json.RawMessage, userID string) {
	budgets := []types.Budget{}

	var requestData struct {
		AccountID string `json:"AccountID"`
	}

	if err := db.DB.Where("user_id =? AND account_id = ?", userID, requestData.AccountID).Find(&budgets).Error; err != nil {
		Message(ws, "Error: Budgets not found")
	}
	var wg sync.WaitGroup
	for i := range budgets {
		wg.Add(1)
		go func(a *types.Budget) {
			defer wg.Done()
			if err := GetBudgetCal(ws, a.ID); err != nil {
				logger.Error("%s", err.Error())
			}
		}(&budgets[i])
	}
	wg.Wait()
	if err := db.DB.Where("user_id = ?", userID).Find(&budgets).Error; err != nil {
		if err := ws.WriteMessage(websocket.TextMessage, []byte(`{"Error":"budgets not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
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
		logger.Error("%s", err.Error())
	}
}

func UpdateBudget(ws *websocket.Conn, data json.RawMessage, userID string) {
	budget := new(types.Budget)

	if err := json.Unmarshal(data, budget); err != nil {
		Message(ws, InvalidData)
	}
	if budget.ID == "" {
		Message(ws, "Error: ID is required")
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, budget.ID).First(&budget).Error; err != nil {
		Message(ws, "Error: Budget not found")
	}

	if err := db.DB.Save(budget).Error; err != nil {
		Message(ws, "Error: Saving")
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
		logger.Error("%s", err.Error())
	}
}

func DeleteBudget(ws *websocket.Conn, data json.RawMessage, userID string) {
	budget := new(types.Budget)

	if err := json.Unmarshal(data, budget); err != nil {
		Message(ws, InvalidData)
	}

	if budget.ID == "" {
		Message(ws, "Error: ID is required")
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, budget.ID).Delete(&budget).Error; err != nil {
		Message(ws, "Error: Failed to delete budget")
	}

	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID": budget.ID,
		},
	}

	responseData, _ := json.Marshal(response)

	if err := ws.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}
}

func GetBudgetCal(ws *websocket.Conn, accountID string) error {
	transactions := []types.Transaction{}
	account := new(types.Accounts)
	budget := new(types.Budget)

	account.ID = accountID

	if err := db.DB.Where(" id =?", account.ID).First(&account).Error; err != nil {
		Message(ws, "Error: Account not found")
		return err
	}

	if err := db.DB.Where("account_id = ?", account.ID).Find(&transactions).Error; err != nil {

		Message(ws, "Error: Could not fetch transactions")
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
		Message(ws, "Error: Failed to save")
		return err
	}

	return nil
}
