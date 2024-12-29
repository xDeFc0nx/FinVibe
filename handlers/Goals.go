package handlers

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/xDeFc0nx/logger-go-pkg"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func CreateGoal(c *websocket.Conn, data json.RawMessage, userID string) {
	goal := new(types.Goal)

	goal.UserID = userID

	var requestData struct {
		AccountID string `json:"AccountID"`
	}
	if err := json.Unmarshal(data, goal); err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid goal data"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

	}
	if goal.ID == "" {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"ID is required"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, requestData.AccountID).First(&goal.Account).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Account not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

	}

	if err := db.DB.Create(&goal).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to create goal"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

	}
	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID":          goal.ID,
			"Amount":      goal.Amount,
			"Description": goal.Description,
		},
	}

	responseData, _ := json.Marshal(response)

	if err := c.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}

}

func GetGoals(c *websocket.Conn, userID string) {
	goals := []types.Goal{}

	var requestData struct {
		AccountID string `json:"AccountID"`
	}

	if err := db.DB.Where("user_id =? AND account_id = ?", userID, requestData.AccountID).Find(&goals).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"goals not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

	}
	var wg sync.WaitGroup
	for i := range goals {
		wg.Add(1)
		go func(a *types.Goal) {
			defer wg.Done()
			if err := GetGoalCal(c, a.ID); err != nil {
				logger.Error("%s", err.Error())
			}
		}(&goals[i])
	}
	wg.Wait()
	if err := db.DB.Where("user_id = ?", userID).Find(&goals).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"goals not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
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
	if err := c.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}

}

func UpdateGoal(c *websocket.Conn, data json.RawMessage, userID string) {

	goal := new(types.Goal)

	if err := json.Unmarshal(data, goal); err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid goal data"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

	}
	if goal.ID == "" {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"ID is required"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, goal.ID).First(&goal).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"goal not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

	}

	if err := db.DB.Save(goal).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"goal not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

	}

	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID":          goal.ID,
			"UserID":      goal.UserID,
			"Goal Amount": goal.GoalAmount,
		},
	}

	responseData, _ := json.Marshal(response)

	if err := c.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}

}

func DeleteGoal(c *websocket.Conn, data json.RawMessage, userID string) {
	goal := new(types.Goal)

	if err := json.Unmarshal(data, goal); err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid goal data"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

	}

	if goal.ID == "" {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"ID is required"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, goal.ID).Delete(&goal).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to delete goal"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

	}

	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID": goal.ID,
		},
	}

	responseData, _ := json.Marshal(response)

	if err := c.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}

}
func GetGoalCal(c *websocket.Conn, accountID string) error {
	transactions := []types.Transaction{}
	account := new(types.Accounts)
	goal := new(types.Goal)

	account.ID = accountID

	if err := db.DB.Where(" id =?", account.ID).First(&account).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Account not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return err
	}

	if err := db.DB.Where("account_id = ?", account.ID).Find(&transactions).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Could not fetch transactions"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
		return err
	}

	totalBalance := float64(0)
	for _, t := range transactions {
		if t.AccountID == accountID {
			totalBalance += t.Amount
		}
	}

	goal.Amount = totalBalance

	if err := db.DB.Save(goal).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"failed to Save"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return err
	}

	return nil
}
