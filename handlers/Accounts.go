package handlers

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func CreateAccount(ws *websocket.Conn, data json.RawMessage, userID string) {
	account := new(types.Accounts)

	user := new(types.User)

	account.ID = uuid.NewString()
	account.UserID = userID

	if err := json.Unmarshal(data, &account); err != nil {
		Send_Error(ws, "Invalid  form data", err)
	}

	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		Send_Error(ws, "User ID Invalid", err)
	}

	if err := db.DB.Create(account).Error; err != nil {
		Send_Error(ws, "Failed to Create Account", err)
	}
	accountData := map[string]interface{}{
		"accountID": account.ID,
		"Type":      account.Type,
	}

	response := map[string]interface{}{
		"account": accountData,
	}

	responseData, _ := json.Marshal(response)
	Send_Message(ws, string(responseData))
}

func GetAccounts(ws *websocket.Conn, data json.RawMessage, userID string) {
	accounts := []types.Accounts{}

	if err := db.DB.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		Send_Error(ws, "Accounts not found", err)
		return
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	for i := range accounts {
		wg.Add(1)
		go func(a *types.Accounts) {
			defer wg.Done()
			if err := GetAccountsBalance(ws, a.ID); err != nil {
				Send_Error(ws, "failed to get account balance", err)
			}
			mu.Lock()
			defer mu.Unlock()
			if err := db.DB.Where("id = ?", a.ID).First(a).Error; err != nil {
				Send_Error(ws, "Failed to fetch updated account", err)
			}
		}(&accounts[i])
	}
	wg.Wait()
	if err := db.DB.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		Send_Error(ws, "Accounts not found", err)
	}

	accountData := make([]map[string]interface{}, len(accounts))

	for i, a := range accounts {
		accountData[i] = map[string]interface{}{
			"ID":             a.ID,
			"UserID":         a.UserID,
			"AccountID":      a.ID,
			"Type":           a.Type,
			"AccountBalance": a.Balance,
			"Income":         a.Income,
			"Expense":        a.Expense,
		}
	}

	// Package the response
	response := map[string]interface{}{
		"accounts": accountData,
	}

	responseData, _ := json.Marshal(response)

	Send_Message(ws, string(responseData))
}

func UpdateAccount(ws *websocket.Conn, data json.RawMessage, userID string) {
	account := new(types.Accounts)
	if err := json.Unmarshal(data, &account); err != nil {
		Send_Error(ws, InvalidData, err)
		return
	}

	if err := db.DB.Where("user_id = ?", userID).Find(&account).Error; err != nil {
		Send_Error(ws, "Account not found", err)
		return
	}
	if err := db.DB.Save(account).Error; err != nil {
		Send_Error(ws, "Failed to update", err)
	}
	accountData := map[string]interface{}{
		"ID":   account.ID,
		"Type": account.Type,
	}

	response := map[string]interface{}{
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

	if err := db.DB.Where("user_id = ?", userID).Find(&account).Error; err != nil {
		Send_Error(ws, "account not found", err)
		return
	}
	if err := db.DB.Delete(account).Error; err != nil {
		Send_Error(ws, "Failed to delete", err)
	}

	accountData := map[string]interface{}{
		"ID":   account.ID,
		"Type": account.Type,
	}

	response := map[string]interface{}{
		"Success": accountData,
	}

	responseData, _ := json.Marshal(response)
	Send_Message(ws, string(responseData))
}
