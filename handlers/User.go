package handlers

import (
	"context"
	"encoding/json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"regexp"
	"time"
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
		return c.Status(400).
			JSON(fiber.Map{"error": "Unable to parse request body", "details": err.Error()})
	}
	emailRegex := regexp.MustCompile(
		`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
	)

	if !emailRegex.MatchString(req.Email) {
		return c.Status(400).
			JSON(fiber.Map{"error": "Invalid email address", "email": req.Email})
	}
	var emailExists bool
	if err := db.DB.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)
		`, req.Email).Scan(&emailExists); err != nil {
		slog.Error("Failed to check email existence", slog.String("error", err.Error()))
	}
	if emailExists {
		data := map[string]any{
			"message": "Email already exists",
		}
		JSendError(c, data, fiber.StatusBadRequest)
	}
	user.ID = uuid.New().String()

	if req.FirstName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "First Name is Required"})
	}

	if req.LastName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Last Name is Required"})
	}
	if req.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Email is Required"})
	}

	if req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Password is required"})
	}
	if req.Currency == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Currency is required"})
	}

	if len(req.Password) < 8 {
		return c.Status(500).
			JSON(fiber.Map{"error": "Password requires at least 8 characters"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Failed to Hash", "details": err.Error()})
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
			"message": "error creating user",
		}
		JSendError(c, data, fiber.StatusBadRequest)

		return nil

	}

	socketID, err := CreateWebSocket(user.ID)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Failed to create WebSocket", "details": err.Error()})
	}
	token, exp, err := Create_JWT_Token(user.ID, socketID)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Failed to Create User", "details": err.Error()})
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
		slog.Error("Failed to fetch user", slog.String("error", err.Error()))
		Send_Error(ws, "failed to get user data", err)
		return
	}

	userData := map[string]interface{}{
		"ID":        user.ID,
		"FirstName": user.FirstName,
		"LastName":  user.LastName,
		"Email":     user.Email,
		"Currency":  user.Currency,
	}

	response := map[string]interface{}{
		"userData": userData,
	}

	responseData, _ := json.Marshal(response)
	Send_Message(ws, string(responseData))
}

func UpdateUser(ws *websocket.Conn, data json.RawMessage, userID string) {
	user := new(types.User)

	if err := db.DB.QueryRow(context.Background(),
		"SELECT * FROM users WHERE id = $1", userID).Scan(&user); err != nil {
		Send_Error(ws, "User not found", err)
		return
	}
	if err := json.Unmarshal(data, &user); err != nil {
		Send_Error(ws, InvalidData, err)
		return
	}

	emailRegex := regexp.MustCompile(
		`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
	)

	if !emailRegex.MatchString(user.Email) {
		Send_Error(ws, "Invalid Email Address", nil)
		return
	}

	var existingUser types.User
	if err := db.DB.QueryRow(context.Background(), "SELECT * FROM users WHERE email = $1 AND id != $2",
		user.Email, userID).Scan(&existingUser); err == nil {
		Send_Error(ws, "Email already exists", nil)
		return
	}
	if user.FirstName == "" {
		Send_Error(ws, "First Name is required", nil)
		return
	}

	if user.LastName == "" {
		Send_Error(ws, "Last Name s required", nil)
		return
	}
	if user.Email == "" {
		Send_Error(ws, "Email is required", nil)
		return
	}

	if user.Currency == "" {
		Send_Error(ws, "Currency is required", nil)
		return
	}

	if len(user.Password) < 8 {
		Send_Error(ws, "Password must be at least 8 characters", nil)
		return
	}

	if _, err := db.DB.Exec(context.Background(),
		"UPDATE users SET firs_tname = $1, last_name = $2, email = $3, password = $4, currency = $5 WHERE id = $6",
		user.FirstName, user.LastName, user.Email, user.Password, user.Currency, userID,
	); err != nil {
		Send_Error(ws, "Failed to save", err)
		return
	}

	userData := map[string]interface{}{
		"ID":        user.ID,
		"FirstName": user.FirstName,
		"LastName":  user.LastName,
		"Email":     user.Email,
	}

	response := map[string]interface{}{
		"Success": userData,
	}

	responseData, _ := json.Marshal(response)
	Send_Message(ws, string(responseData))
}

func DeleteUser(ws *websocket.Conn, data json.RawMessage, userID string) {
	user := new(types.User)

	if err := db.DB.QueryRow(context.Background(),
		"SELECT * FROM users WHERE id = $1", userID).Scan(&user); err != nil {
		Send_Error(ws, "User not found", err)
		return
	}
	if _, err := db.DB.Exec(context.Background(), "DELETE users WHERE id = $1", userID); err != nil {
		Send_Error(ws, "Failed to update user", err)
		return
	}
	response := map[string]string{
		"Success": "Deleted User",
	}
	responseJson, _ := json.Marshal(response)
	Send_Message(ws, string(responseJson))
}
