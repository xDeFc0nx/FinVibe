package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func CreateGoal(ws *websocket.Conn, data json.RawMessage, userID string) {
	goal := new(types.Goal)

	goal.UserID = userID

	type requestData struct {
		AccountID string `json:"AccountID"`
	}

	var req requestData
	if err := json.Unmarshal(data, goal); err != nil {
		Send_Error(ws, InvalidData, err)
	}
	if goal.ID == "" {
		Send_Error(ws, MsgMissingID, nil)
	}
	if _, err := db.DB.Exec(context.Background(), `
 SELECT FROM goals WHERE id = $1 AND user_id = $2
		`, req.AccountID, userID); err != nil {
		Send_Error(ws, MsgAccountNotFound, err)
	}
	if _, err := db.DB.Exec(context.Background(), `
INSERT INTO goals (id, user_id, account_id, amount, description)
		VALUES ($1, $2, $3, $4, $5) 
		`,
		uuid.New().String(),
		goal.UserID,
		goal.AccountID,
		goal.Amount,
		goal.Description,
	); err != nil {
		Send_Error(ws, fmt.Sprintf(MsgCreateFailedFmt, "Goal"), err)
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

	type requestData struct {
		AccountID string `json:"AccountID"`
	}
	var req requestData
	if _, err := db.DB.Exec(context.Background(), `
SELECT FROM accounts WHERE id = $1 AND user_id = $2
		`, req.AccountID, userID); err != nil {
		Send_Error(ws, MsgAccountNotFound, err)
	}
	var wg sync.WaitGroup
	for i := range goals {
		wg.Add(1)
		go func(a *types.Goal) {
			defer wg.Done()
			if err := GetGoalCal(ws, a.ID); err != nil {
				Send_Error(ws, fmt.Sprintf(MsgFetchFailedFmt, "Goal"), err)
			}
		}(&goals[i])
	}
	wg.Wait()
	if _, err := db.DB.Exec(context.Background(), `
SELECT FROM goals WHERE user_id = $1
		`, userID); err != nil {
		Send_Error(ws, fmt.Sprintf(MsgFetchFailedFmt, "Goal"), err)
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
	type requestData struct {
		AccountID  string  `json:"AccountID"`
		Ammount    float64 `json:"Ammount"`
		GoalAmount float64 `json:"GoalAmount"`
	}
	var req requestData

	goal := new(types.Goal)

	if err := json.Unmarshal(data, goal); err != nil {
		Send_Error(ws, MsgInvalidData, err)
	}
	if goal.ID == "" {
		Send_Error(ws, fmt.Sprintf(MsgMissingFieldFmt, "Goal ID"), nil)
	}
	if _, err := db.DB.Exec(context.Background(), `
SELECT FROM goals WHERE id = $1 AND user_id = $2
		`, goal.ID, userID); err != nil {
		Send_Error(ws, "Goal not found", err)
	}
	if _, err := db.DB.Exec(context.Background(), `
INSERT INTO goals (id, user_id, account_id, amount, description)
VALUES ($1, $2, $3, $4, $5)
		`,
		goal.ID,
		userID,
		req.AccountID,
		req.Ammount,
		req.GoalAmount,
	); err != nil {
		Send_Error(ws, "failed to update goal", err)
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
		Send_Error(ws, MsgInvalidData, err)
	}

	if goal.ID == "" {
		Send_Error(ws, "ID is required", nil)
	}
	if _, err := db.DB.Exec(context.Background(), `
FROM goals WHERE id = $1 AND user_id = $2
		`, goal.ID, userID); err != nil {
		Send_Error(ws, MsgBudgetNotFound, err)
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
	if _, err := db.DB.Exec(context.Background(), `
SELECT 1 FROM accounts WHERE id = $1
		`, accountID); err != nil {
		Send_Error(ws, MsgAccountNotFound, err)
	}
	rows, err := db.DB.Query(context.Background(),
		`SELECT id, user_id, account_id, amount, description, is_recurring,
		created_at, FROM transactions WHERE account_id = $1`, account.ID,
	)

	if err != nil {
		Send_Error(ws, MsgAccountNotFound, err)
	}

	defer rows.Close()
	transactions, err = pgx.CollectRows(rows,
		pgx.RowTo[types.Transaction])
	if err != nil {
	}
	totalBalance := float64(0)
	for _, t := range transactions {
		if t.AccountID == accountID {
			totalBalance += t.Amount
		}
	}

	goal.Amount = totalBalance

	if _, err := db.DB.Exec(context.Background(), `
UPDATE goals SET amount = $1 WHERE account_id = $2
		`, goal.Amount, accountID); err != nil {
		Send_Error(ws, fmt.Sprintf(MsgUpdateFailedFmt, "Goal"), err)
	}

	return nil
}
