package handlers

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func getCharts(ws *websocket.Conn, data json.RawMessage, userID string) {
	if err := json.Unmarshal(data, &requestData); err != nil {
		Send_Error(ws, "Invalid data format", err)
		return
	}

	account := new(types.Accounts)
	if err := db.DB.Where("user_id = ? AND id = ?", userID, requestData.AccountID).First(account).Error; err != nil {
		Send_Error(ws, "Account not found", err)
		return
	}

	start, end := GetDateRange(requestData.DateRange)

	var transactions []types.Transaction
	if err := db.DB.Where("account_id = ? AND created_at BETWEEN ? AND ?",
		requestData.AccountID,
		start,
		end,
	).Find(&transactions).Error; err != nil {
		Send_Error(ws, "Failed to retrieve transactions", err)
		return
	}

	type byDay struct {
		Day     int     `json:"day"`
		Income  float64 `json:"Income"`
		Expense float64 `json:"Expenses"`
	}

	var formattedData []byDay

	dailyTotals := make(map[int]*byDay)

	for _, t := range transactions {
		day := t.CreatedAt.Day()
		if dailyTotals[day] == nil {
			dailyTotals[day] = &byDay{
				Day:     day,
				Income:  0.0,
				Expense: 0.0,
			}
		}
		switch t.Type {
		case "Income":
			dailyTotals[day].Income += t.Amount
		case "Expense":
			dailyTotals[day].Expense += t.Amount
		}
	}
	for _, totals := range dailyTotals {
		formattedData = append(formattedData, *totals)
	}

	response := map[string]interface{}{
		"chartData": formattedData,
	}

	responseData, _ := json.Marshal(response)

	Send_Message(ws, string(responseData))
}
