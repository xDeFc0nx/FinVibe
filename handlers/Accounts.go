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

func CreateAccount(ws *websocket.Conn, data json.RawMessage, userID string) {
	account := new(types.Accounts)

	user := new(types.User)

	account.ID = uuid.NewString()
	account.UserID = userID

	if err := json.Unmarshal(data, &account); err != nil {
		Message(ws, "Invalid  form data", err)
	}

	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		Message(ws, "User ID Invalid", err)
	}

	if err := db.DB.Create(account).Error; err != nil {
		Message(ws, "Failed to Create Account", err)
	}
	response, _ := json.Marshal(account)
	if err := ws.WriteMessage(websocket.TextMessage, response); err != nil {
		slog.Error("Failed to send message", slog.String("error", err.Error()))
	}
}

func GetAccounts(ws *websocket.Conn, data json.RawMessage, userID string) {
	accounts := []types.Accounts{}

	if err := db.DB.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		Message(ws, "Accounts not found", err)
		return
	}

	var wg sync.WaitGroup
	for i := range accounts {
		wg.Add(1)
		go func(a *types.Accounts) {
			defer wg.Done()
			if err := GetAccountBalance(ws, a.ID); err != nil {
				Message(ws, "failed to get account balance", err)
			}
		}(&accounts[i])
	}
	wg.Wait()
	if err := db.DB.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		Message(ws, "Accounts not found", err)
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
		slog.Error("Failed to send message", slog.String("error", err.Error()))
	}
}

func UpdateAccount(ws *websocket.Conn, data json.RawMessage, userID string) {
	account := new(types.Accounts)
	if err := json.Unmarshal(data, &account); err != nil {
		Message(ws, InvalidData, err)
		return
	}

	if err := db.DB.Where("user_id = ?", userID).Find(&account).Error; err != nil {
		Message(ws, "Account not found", err)
		return
	}
	if err := db.DB.Save(account).Error; err != nil {
		Message(ws, "Failed to update", err)
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
		slog.Error("Failed to send message", slog.String("error", err.Error()))
	}
}

func DeleteAccount(ws *websocket.Conn, data json.RawMessage, userID string) {
	account := new(types.Accounts)
	if err := json.Unmarshal(data, &account); err != nil {
		Message(ws, InvalidData, err)

		return

	}

	if err := db.DB.Where("user_id = ?", userID).Find(&account).Error; err != nil {
		Message(ws, "account not found", err)
		return
	}
	if err := db.DB.Delete(account).Error; err != nil {
		if err := ws.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to Delete"}`+err.Error())); err != nil {
			slog.Error(
				"Failed to send message",
				slog.String("error", err.Error()),
			)
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
		slog.Error("Failed to send message", slog.String("error", err.Error()))
	}
}
