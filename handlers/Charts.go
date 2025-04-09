package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func getCharts(ws *websocket.Conn, data json.RawMessage, userID string) {
	if err := json.Unmarshal(data, &requestData); err != nil {
		Send_Error(ws, MsgInvalidData, err)
		return
	}

	if _, err := db.DB.Exec(context.Background(), `
SELECT EXISTS (SELECT 1 FROM accounts WHERE id = $1 AND user_id = $2)
		`, requestData.AccountID, userID); err != nil {
		Send_Error(ws, MsgAccountNotFound, err)
	}
	start, end := GetDateRange(requestData.DateRange)
	rows, err := db.DB.Query(context.Background(), `
SELECT amount, id, user_id, account_id, type, description, is_recurring, created_at, updated_at
		FROM transactions
		WHERE account_id = $1 
		AND created_at BETWEEN $2 AND $3
		ORDER BY created_at DESC`,
		requestData.AccountID,
		start,
		end,
	)
	if err != nil {
		Send_Error(ws, fmt.Sprintf(MsgFetchFailedFmt, "transactions"), err)
	}

	defer rows.Close()
	var transactions []types.Transaction
	transactions, err = pgx.CollectRows(rows, pgx.RowToStructByName[types.Transaction])
	if err != nil {
		Send_Error(ws, MsgCollectRowsFailed, err)
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
	type byDesc struct {
		Description string  `json:"Description"`
		Amount      float64 `json:"Amount"`
	}

	var Income []byDesc
	var Expenses []byDesc

	IncomeTotals := make(map[string]float64)
	ExpensesTotals := make(map[string]float64)
	for _, t := range transactions {
		switch t.Type {
		case "Income":
			IncomeTotals[t.Description] += t.Amount
		case "Expense":
			ExpensesTotals[t.Description] += t.Amount
		}
	}
	for desc, amount := range IncomeTotals {
		Income = append(Income, byDesc{
			Description: desc,
			Amount:      amount,
		})
	}
	for desc, amount := range ExpensesTotals {
		Expenses = append(Expenses, byDesc{
			Description: desc,
			Amount:      amount,
		})
	}
	response := map[string]interface{}{
		"chartData":   formattedData,
		"IncomePie":   Income,
		"ExpensesPie": Expenses,
	}

	responseData, _ := json.Marshal(response)

	Send_Message(ws, string(responseData))
}
