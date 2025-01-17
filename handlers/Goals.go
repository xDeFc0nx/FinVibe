package handlers

import (
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/gofiber/contrib/websocket"

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
		Send_Error(ws, InvalidData, err)
	}
	if goal.ID == "" {
		Send_Error(ws, "ID is Required", nil)
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, requestData.AccountID).First(&goal.Account).Error; err != nil {
		Send_Error(ws, "Account not found", err)
	}

	if err := db.DB.Create(&goal).Error; err != nil {
		Send_Error(ws, "Failed To Create Goal", err)
	}
	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID":          goal.ID,
			"Amount":      goal.Amount,
			"Description": goal.Description,
		},
	}

	responseData, _ := json.Marshal(response)

	Send_Message(ws, string(responseData))
}

func GetGoals(ws *websocket.Conn, data json.RawMessage, userID string) {
	goals := []types.Goal{}

	var requestData struct {
		AccountID string `json:"AccountID"`
	}

	if err := db.DB.Where("user_id =? AND account_id = ?", userID, requestData.AccountID).Find(&goals).Error; err != nil {
		Send_Error(ws, "Goals Not Found", err)
	}
	var wg sync.WaitGroup
	for i := range goals {
		wg.Add(1)
		go func(a *types.Goal) {
			defer wg.Done()
			if err := GetGoalCal(ws, a.ID); err != nil {
				slog.Error(
					"failed to get goal",
					slog.String("error", err.Error()),
				)
			}
		}(&goals[i])
	}
	wg.Wait()
	if err := db.DB.Where("user_id = ?", userID).Find(&goals).Error; err != nil {
		Send_Error(ws, "Goals not found", err)
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

	response := map[string]interface{}{
		"Success": "Fetched goals",
		"goals":   goalsData,
	}

	responseData, _ := json.Marshal(response)
	Send_Message(ws, string(responseData))
}

func UpdateGoal(ws *websocket.Conn, data json.RawMessage, userID string) {
	goal := new(types.Goal)

	if err := json.Unmarshal(data, goal); err != nil {
		Send_Error(ws, "Invalid goal data", err)
	}
	if goal.ID == "" {
		Send_Error(ws, "Goal is required", nil)
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, goal.ID).First(&goal).Error; err != nil {
		Send_Error(ws, "Goal not found", err)
	}

	if err := db.DB.Save(goal).Error; err != nil {
		Send_Error(ws, "Goal not found", err)
	}

	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID":          goal.ID,
			"UserID":      goal.UserID,
			"Goal Amount": goal.GoalAmount,
		},
	}

	responseData, _ := json.Marshal(response)

	Send_Message(ws, string(responseData))
}

func DeleteGoal(ws *websocket.Conn, data json.RawMessage, userID string) {
	goal := new(types.Goal)

	if err := json.Unmarshal(data, goal); err != nil {
		Send_Error(ws, "Invalid goal data", err)
	}

	if goal.ID == "" {
		Send_Error(ws, "ID is required", nil)
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, goal.ID).Delete(&goal).Error; err != nil {
		Send_Error(ws, "Failed to delete goal", err)
	}

	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID": goal.ID,
		},
	}

	responseData, _ := json.Marshal(response)
	Send_Message(ws, string(responseData))
}

func GetGoalCal(ws *websocket.Conn, accountID string) error {
	transactions := []types.Transaction{}
	account := new(types.Accounts)
	goal := new(types.Goal)

	account.ID = accountID

	if err := db.DB.Where(" id =?", account.ID).First(&account).Error; err != nil {
		Send_Error(ws, "Account Not found", err)
	}

	if err := db.DB.Where("account_id = ?", account.ID).Find(&transactions).Error; err != nil {
		Send_Error(ws, "Could not fetch transactions", err)
	}

	totalBalance := float64(0)
	for _, t := range transactions {
		if t.AccountID == accountID {
			totalBalance += t.Amount
		}
	}

	goal.Amount = totalBalance

	if err := db.DB.Save(goal).Error; err != nil {
		Send_Error(ws, "Failed to save", err)
	}
	return nil
}
