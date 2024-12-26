package handlers

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/xDeFc0nx/logger-go-pkg"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func CreateTransaction(c *websocket.Conn, data json.RawMessage, userID string) {

	transaction := new(types.Transaction)
	account := new(types.Accounts)
	transaction.ID = uuid.NewString()

	transaction.UserID = userID

	if err := json.Unmarshal(data, &transaction); err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid transaction data"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	if err := db.DB.Where("user_id =? AND account_id =?", userID, account.ID).First(&account).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Account not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}
	if err := db.DB.Create(&transaction).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to create transaction"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	response, _ := json.Marshal(transaction)
	if err := c.WriteMessage(websocket.TextMessage, response); err != nil {
		logger.Error("%s", err.Error())
	}

}
func GetTransactions(c *websocket.Conn, data json.RawMessage, userID string) {
	transactions := []types.Transaction{}
	account := new(types.Accounts)

	if err := json.Unmarshal(data, &transactions); err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid transaction data"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}
	if err := db.DB.Where("user_id =? AND account_id =?", userID, account.ID).First(&account).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Account not found"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	response, _ := json.Marshal(transactions)
	if err := c.WriteMessage(websocket.TextMessage, response); err != nil {
		logger.Error("%s", err.Error())
	}
}

func GetTransactionById(c *websocket.Conn, data json.RawMessage, userID string) {
	transaction := new(types.Transaction)

	transaction.UserID = userID

	if err := json.Unmarshal(data, &transaction); err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid transaction data"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}
	if transaction.UserID != userID {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction does not belong to the user"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	if err := db.DB.Where("id = ?", transaction.ID).First(&transaction).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction not found"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	response, _ := json.Marshal(transaction)
	if err := c.WriteMessage(websocket.TextMessage, response); err != nil {
		logger.Error("%s", err.Error())
	}
}

func UpdateTransction(c *websocket.Conn, data json.RawMessage, userID string) {
	transaction := new(types.Transaction)

	if err := json.Unmarshal(data, &transaction); err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error": "Invalid Data"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
	}

	if transaction.UserID != userID {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction does not belong to the user"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	if err := db.DB.Where("id = ?", transaction.ID).First(&transaction).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return

	}

	if err := db.DB.Save(transaction).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Failed to Update"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
	}
	if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Success":"Transaction Updated"}`)); err != nil {
		logger.Error("%s", err.Error())
	}
}
func DeleteTransaction(c *websocket.Conn, data json.RawMessage, userID string) {
	transaction := new(types.Transaction)

	transaction.UserID = userID

	if err := json.Unmarshal(data, &transaction); err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Invalid transaction data"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}

		return
	}
	if transaction.UserID != userID {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction does not belong to the user"}`)); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	if err := db.DB.Where("id = ?", transaction.ID).First(&transaction).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"Transaction not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	if err := db.DB.Delete(transaction).Error; err != nil {
		if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Error":"failed to Delete"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	if err := c.WriteMessage(websocket.TextMessage, []byte(`{"Success":"Transaction Deleted"}`)); err != nil {
		logger.Error("%s", err.Error())
	}

}
