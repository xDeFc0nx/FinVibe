package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *fiber.Ctx) error {
	type Request struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Currency  string `json:"currency"`
	}
	req := Request{}

	user := new(types.User)
	ctx := c.Context()
	if err := c.BodyParser(&req); err != nil {
		data := map[string]any{
			"message": MsgInvalidData,
		}
		JSendError(c, data, fiber.StatusBadRequest, err)
		return nil
	}
	emailRegex := regexp.MustCompile(
		`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
	)

	if !emailRegex.MatchString(req.Email) {
		data := map[string]any{
			"message": MsgInvalidEmail,
		}
		JSendError(c, data, fiber.StatusBadRequest, nil)
		return nil

	}
	var emailExists bool
	if err := db.DB.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)
		`, req.Email).Scan(&emailExists); err != nil {
		data := map[string]any{
			"message": "Failed to check email existence",
		}
		JSendError(c, data, fiber.StatusBadRequest, err)
	}
	if emailExists {
		data := map[string]any{
			"message": MsgEmailExists,
		}
		JSendError(c, data, fiber.StatusBadRequest, nil)
		return nil
	}
	user.ID = uuid.New().String()

	if req.FirstName == "" {
		data := map[string]any{
			"message": fmt.Sprintf(MsgMissingFieldFmt, "First Name"),
		}
		JSendError(c, data, fiber.StatusBadRequest, nil)
		return nil
	}

	if req.LastName == "" {
		data := map[string]any{
			"message": fmt.Sprintf(MsgMissingFieldFmt, "Last Name"),
		}
		JSendError(c, data, fiber.StatusBadRequest, nil)
		return nil

	}
	if req.Email == "" {
		data := map[string]any{
			"message": fmt.Sprintf(MsgMissingFieldFmt, "Email"),
		}
		JSendError(c, data, fiber.StatusBadRequest, nil)
		return nil

	}

	if req.Password == "" {

		data := map[string]any{
			"message": fmt.Sprintf(MsgMissingFieldFmt, "Password"),
		}
		JSendError(c, data, fiber.StatusBadRequest, nil)
		return nil
	}
	if req.Currency == "" {
		data := map[string]any{
			"message": fmt.Sprintf(MsgMissingFieldFmt, "Currency"),
		}
		JSendError(c, data, fiber.StatusBadRequest, nil)
		return nil
	}

	if len(req.Password) < 8 {
		data := map[string]any{
			"message": fmt.Sprintf(MsgMissingFieldFmt, "Password must be 8 characters"),
		}
		JSendError(c, data, fiber.StatusBadRequest, nil)
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		data := map[string]any{
			"message": "Failed to hash",
		}
		JSendError(c, data, fiber.StatusBadRequest, err)
		return nil
	}

	if _, err := db.DB.Exec(context.Background(),
		"INSERT INTO users (id, first_name, last_name, email, password, currency, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.ID,
		req.FirstName,
		req.LastName,
		req.Email,
		hashedPassword,
		req.Currency,
		time.Now().UTC(),
	); err != nil {
		data := map[string]any{
			"message": fmt.Sprintf(MsgCreateFailedFmt, "User"),
		}
		JSendError(c, data, fiber.StatusInternalServerError, err)

		return nil

	}

	socketID, err := CreateWebSocket(user.ID)
	if err != nil {
		data := map[string]any{
			"message": fmt.Sprintf(MsgCreateFailedFmt, "WebSocket"),
		}
		JSendError(c, data, fiber.StatusInternalServerError, err)
		return nil
	}
	token, exp, err := CreateJWTToken(user.ID, socketID)
	if err != nil {
		data := map[string]any{
			"message": fmt.Sprintf(MsgCreateFailedFmt, "JWT Token"),
		}
		JSendError(c, data, fiber.StatusInternalServerError, err)
		return nil
	}
	slog.Info(
		"User Created",
		"userID",
		user.ID,
	)
	return c.JSON(fiber.Map{
		"message":  "User created successfully",
		"userID":   user.ID,
		"socketID": socketID,
		"token":    token,
		"exp":      exp,
	})
}

func GetUser(ws *websocket.Conn, data json.RawMessage, userID string) {
	user := new(types.User)
	if err := db.DB.QueryRow(context.Background(), `
    SELECT id, first_name, last_name, email, password, currency, created_at 
    FROM users 
    WHERE id = $1
`, userID).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Currency, &user.CreatedAt); err != nil {
		Send_Error(ws, fmt.Sprintf(MsgFetchFailedFmt, "User"), err)
		return
	}

	userData := map[string]any{
		"ID":        user.ID,
		"FirstName": user.FirstName,
		"LastName":  user.LastName,
		"Email":     user.Email,
		"Currency":  user.Currency,
	}

	response := map[string]any{
		"userData": userData,
	}

	responseData, _ := json.Marshal(response)
	Send_Message(ws, string(responseData))
}

func UpdateUser(ws *websocket.Conn, data json.RawMessage, userID string) {
	user := new(types.User)

	if err := db.DB.QueryRow(context.Background(),
		"SELECT * FROM users WHERE id = $1", userID).Scan(&user); err != nil {
		Send_Error(ws, MsgUserNotFound, err)
		return
	}
	if err := json.Unmarshal(data, &user); err != nil {
		Send_Error(ws, MsgInvalidData, err)
		return
	}

	emailRegex := regexp.MustCompile(
		`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
	)

	if !emailRegex.MatchString(user.Email) {
		Send_Error(ws, MsgInvalidEmail, nil)
		return
	}

	var existingUser types.User
	if err := db.DB.QueryRow(context.Background(), "SELECT * FROM users WHERE email = $1 AND id != $2",
		user.Email, userID).Scan(&existingUser); err == nil {
		Send_Error(ws, MsgEmailExists, nil)
		return
	}
	if user.FirstName == "" {
		Send_Error(ws, fmt.Sprintf(MsgMissingFieldFmt, "First Name"), nil)
		return
	}

	if user.LastName == "" {
		Send_Error(ws, fmt.Sprintf(MsgMissingFieldFmt, "Last Name"), nil)
		return
	}
	if user.Email == "" {
		Send_Error(ws, fmt.Sprintf(MsgMissingFieldFmt, "Email"), nil)
		return
	}

	if user.Currency == "" {
		Send_Error(ws, fmt.Sprintf(MsgMissingFieldFmt, "Currency"), nil)
		return
	}

	if len(user.Password) < 8 {
		Send_Error(ws, fmt.Sprintf(MsgMissingFieldFmt, "Password must be at least 8 characters"), nil)
		return
	}

	if _, err := db.DB.Exec(context.Background(),
		"UPDATE users SET firs_tname = $1, last_name = $2, email = $3, password = $4, currency = $5 WHERE id = $6",
		user.FirstName, user.LastName, user.Email, user.Password, user.Currency, userID,
	); err != nil {
		Send_Error(ws, fmt.Sprintf(MsgUpdateFailedFmt, "User"), err)
		return
	}

	userData := map[string]any{
		"ID":        user.ID,
		"FirstName": user.FirstName,
		"LastName":  user.LastName,
		"Email":     user.Email,
	}

	response := map[string]any{
		"Success": userData,
	}

	responseData, _ := json.Marshal(response)
	Send_Message(ws, string(responseData))
}

func DeleteUser(ws *websocket.Conn, data json.RawMessage, userID string) {
	user := new(types.User)

	if err := db.DB.QueryRow(context.Background(),
		"SELECT * FROM users WHERE id = $1", userID).Scan(&user); err != nil {
		Send_Error(ws, MsgUserNotFound, err)
		return
	}
	if _, err := db.DB.Exec(context.Background(), "DELETE users WHERE id = $1", userID); err != nil {
		Send_Error(ws, fmt.Sprintf(MsgDeleteFailedFmt, "User"), err)
		return
	}
	response := map[string]string{
		"Success": "Deleted User",
	}
	responseJson, _ := json.Marshal(response)
	Send_Message(ws, string(responseJson))
}
