package handlers

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/xDeFc0nx/logger-go-pkg"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func CreateGoal(ws *websocket.Conn, data json.RawMessage, userID string) {
	goal := new(types.Goal)

	goal.UserID = userID

	var requestData struct {
		AccountID string `json:"AccountID"`
	}
	if err := json.Unmarshal(data, goal); err != nil {
		Message(ws, InvalidData)
	}
	if goal.ID == "" {
		Message(ws, "Error: ID is Required")
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, requestData.AccountID).First(&goal.Account).Error; err != nil {
		Message(ws, "Error: Account not found")
	}

	if err := db.DB.Create(&goal).Error; err != nil {
		Message(ws, "Error: Failed To Create Goal")
	}
	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID":          goal.ID,
			"Amount":      goal.Amount,
			"Description": goal.Description,
		},
	}

	responseData, _ := json.Marshal(response)

	if err := ws.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}
}

func GetGoals(ws *websocket.Conn, data json.RawMessage, userID string) {
	goals := []types.Goal{}

	var requestData struct {
		AccountID string `json:"AccountID"`
	}

	if err := db.DB.Where("user_id =? AND account_id = ?", userID, requestData.AccountID).Find(&goals).Error; err != nil {
		Message(ws, "Error: Goals Not Found")
	}
	var wg sync.WaitGroup
	for i := range goals {
		wg.Add(1)
		go func(a *types.Goal) {
			defer wg.Done()
			if err := GetGoalCal(ws, a.ID); err != nil {
				logger.Error("%s", err.Error())
			}
		}(&goals[i])
	}
	wg.Wait()
	if err := db.DB.Where("user_id = ?", userID).Find(&goals).Error; err != nil {
		if err := ws.WriteMessage(websocket.TextMessage, []byte(`{"Error":"goals not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		Message(ws, "Error: Goals not found")
	}

	goalsData := make([]map[string]interface{}, len(goals))

	for i, a := range goals {
		goalsData[i] = map[string]interface{}{
			"ID":          a.ID,
			"UserID":      a.UserID,
			"AccountID":   a.ID,
			"Goal Amount": float64(a.Amount),
			"goal":        float64(a.GoalAmount),
		}
	}

	// Package the response
	response := map[string]interface{}{
		"Success": "Fetched goals",
		"goals":   goalsData,
	}

	responseData, _ := json.Marshal(response)
	if err := ws.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}
}

func UpdateGoal(ws *websocket.Conn, data json.RawMessage, userID string) {
	goal := new(types.Goal)

	if err := json.Unmarshal(data, goal); err != nil {
		Message(ws, "Error: Invalid goal data")
	}
	if goal.ID == "" {
		Message(ws, "Error: ID is required")
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, goal.ID).First(&goal).Error; err != nil {
		if err := ws.WriteMessage(websocket.TextMessage, []byte(`{"Error":"goal not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		Message(ws, "Error: Goal not found")
	}

	if err := db.DB.Save(goal).Error; err != nil {
		Message(ws, "Error: Goal not found")
	}

	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID":          goal.ID,
			"UserID":      goal.UserID,
			"Goal Amount": goal.GoalAmount,
		},
	}

	responseData, _ := json.Marshal(response)

	if err := ws.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}
}

func DeleteGoal(ws *websocket.Conn, data json.RawMessage, userID string) {
	goal := new(types.Goal)

	if err := json.Unmarshal(data, goal); err != nil {
		Message(ws, "Error: Invalid goal data")
	}

	if goal.ID == "" {
		Message(ws, "Error: ID is required")
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, goal.ID).Delete(&goal).Error; err != nil {
		Message(ws, "Error: Failed to delete goal")
	}

	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID": goal.ID,
		},
	}

	responseData, _ := json.Marshal(response)

	if err := ws.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}
}

func GetGoalCal(ws *websocket.Conn, accountID string) error {
	transactions := []types.Transaction{}
	account := new(types.Accounts)
	goal := new(types.Goal)

	account.ID = accountID

	if err := db.DB.Where(" id =?", account.ID).First(&account).Error; err != nil {
		Message(ws, "Error: Account Not found")
	}

	if err := db.DB.Where("account_id = ?", account.ID).Find(&transactions).Error; err != nil {
		Message(ws, "Error: Could not fetch transactions")
	}

	totalBalance := float64(0)
	for _, t := range transactions {
		if t.AccountID == accountID {
			totalBalance += t.Amount
		}
	}

	goal.Amount = totalBalance

	if err := db.DB.Save(goal).Error; err != nil {
		Message(ws, "Error: Failed to save")
	}
	return nil
}
