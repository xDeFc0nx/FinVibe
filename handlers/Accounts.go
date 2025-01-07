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

func CreateAccount(ws *websocket.Conn, data json.RawMessage, userID string) {
	account := new(types.Accounts)

	user := new(types.User)

	account.ID = uuid.NewString()
	account.UserID = userID

	if err := json.Unmarshal(data, &account); err != nil {
		Message(ws, "Error: Invalid  form data")
	}

	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		Message(ws, "Error: User ID Invalid")
	}

	if err := db.DB.Create(account).Error; err != nil {
		Message(ws, "Error: Failed to Create Account")
	}
	response, _ := json.Marshal(account)
	if err := ws.WriteMessage(websocket.TextMessage, response); err != nil {
		logger.Error("%s", err.Error())
	}
}

func GetAccounts(ws *websocket.Conn, data json.RawMessage, userID string) {
	accounts := []types.Accounts{}

	if err := db.DB.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		Message(ws, "Error: Accounts not found")
		return
	}

	var wg sync.WaitGroup
	for i := range accounts {
		wg.Add(1)
		go func(a *types.Accounts) {
			defer wg.Done()
			if err := GetAccountBalance(ws, a.ID); err != nil {
				logger.Error("%s", err.Error())
			}
		}(&accounts[i])
	}
	wg.Wait()
	if err := db.DB.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		Message(ws, "Error: Accounts not found")
	}

	accountData := make([]map[string]interface{}, len(accounts))

	for i, a := range accounts {
		accountData[i] = map[string]interface{}{
			"ID":             a.ID,
			"UserID":         a.UserID,
			"AccountID":      a.ID,
			"AccountBalance": float64(a.Balance),
		}
	}

	// Package the response
	response := map[string]interface{}{
		"Success":  "Fetched Accounts",
		"Accounts": accountData,
	}

	responseData, _ := json.Marshal(response)
	if err := ws.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}
}

func UpdateAccount(ws *websocket.Conn, data json.RawMessage, userID string) {
	account := new(types.Accounts)
	if err := json.Unmarshal(data, &account); err != nil {
		Message(ws, "Error: Invalid form data")
		return
	}

	if err := db.DB.Where("user_id = ?", userID).Find(&account).Error; err != nil {
		Message(ws, "Error: Account not found")
		return
	}
	if err := db.DB.Save(account).Error; err != nil {
		Message(ws, "Error: Failed to update")
	}
	accountData := map[string]interface{}{
		"ID":   account.ID,
		"Type": account.Type,
	}

	response := map[string]interface{}{
		"Success": accountData,
	}

	responseData, _ := json.Marshal(response)
	if err := ws.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}
}

func DeleteAccount(ws *websocket.Conn, data json.RawMessage, userID string) {
	account := new(types.Accounts)
	if err := json.Unmarshal(data, &account); err != nil {
		Message(ws, "Error: Invalid form data")

		return

	}

	if err := db.DB.Where("user_id = ?", userID).Find(&account).Error; err != nil {
		if err := ws.WriteMessage(websocket.TextMessage, []byte(`{"Error":"account not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		Message(ws, "")
		return
	}
	if err := db.DB.Delete(account).Error; err != nil {
		if err := ws.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to Delete"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
	}

	accountData := map[string]interface{}{
		"ID":   account.ID,
		"Type": account.Type,
	}

	response := map[string]interface{}{
		"Success": accountData,
	}

	responseData, _ := json.Marshal(response)
	if err := ws.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}
}
