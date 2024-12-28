package handlers

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/xDeFc0nx/logger-go-pkg"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func CreateBudget(c *websocket.Conn, data json.RawMessage, userID string) {

	budget := new(types.Budget)

	budget.ID = uuid.NewString()
	budget.UserID = userID

	var requestData struct {
		AccountID string `json:"AccountID"`
	}

	if err := json.Unmarshal(data, budget); err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid budget data"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

	}

	if budget.ID == "" {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"ID is required"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, requestData.AccountID).First(&budget.Account).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Account not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

	}

	if err := db.DB.Create(&budget).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to create budget"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

	}
	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID":          budget.ID,
			"Amount":      budget.Amount,
			"Description": budget.Description,
		},
	}

	responseData, _ := json.Marshal(response)

	if err := c.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}

}

func GetBudgets(c *websocket.Conn, userID string) {
	budgets := []types.Budget{}

	var requestData struct {
		AccountID string `json:"AccountID"`
	}

	if err := db.DB.Where("user_id =? AND account_id = ?", userID, requestData.AccountID).Find(&budgets).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Budgets not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

	}

	response := map[string]interface{}{
		"Budgets": budgets,
	}

	responseData, _ := json.Marshal(response)

	if err := c.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}

}

func UpdateBudget(c *websocket.Conn, data json.RawMessage, userID string) {

	budget := new(types.Budget)

	if err := json.Unmarshal(data, budget); err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid budget data"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

	}
	if budget.ID == "" {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"ID is required"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, budget.ID).First(&budget).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Budget not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

	}

	if err := db.DB.Save(budget).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Budget not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

	}

	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID":     budget.ID,
			"UserID": budget.UserID,
			"Amount": budget.Amount,
		},
	}

	responseData, _ := json.Marshal(response)

	if err := c.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}

}

func DeleteBudget(c *websocket.Conn, data json.RawMessage, userID string) {
	budget := new(types.Budget)

	if err := json.Unmarshal(data, budget); err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid budget data"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

	}

	if budget.ID == "" {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"ID is required"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
	}

	if err := db.DB.Where("user_id =? AND id =?", userID, budget.ID).Delete(&budget).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to delete budget"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

	}

	response := map[string]interface{}{
		"Success": map[string]interface{}{
			"ID": budget.ID,
		},
	}

	responseData, _ := json.Marshal(response)

	if err := c.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}

}
