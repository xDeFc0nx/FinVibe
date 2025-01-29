package handlers

import (
	"encoding/json"
	"regexp"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func CreateUser(c *fiber.Ctx) error {
	user := new(types.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).
			JSON(fiber.Map{"error": "Unable to parse request body", "details": err.Error()})
	}
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()

	emailRegex := regexp.MustCompile(
		`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
	)

	if !emailRegex.MatchString(user.Email) {
		return c.Status(400).
			JSON(fiber.Map{"error": "Invalid email address", "email": user.Email})
	}

	if err := db.DB.Where(user.Email).Error; err != nil {
		return c.Status(409).JSON(fiber.Map{"error": "Email already exists"})
	}

	if user.FirstName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "First Name is Required"})
	}

	if user.LastName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Last Name is Required"})
	}
	if user.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Email is Required"})
	}

	if user.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Password is required"})
	}
	if user.Currency == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Currency is required"})
	}

	if len(user.Password) < 8 {
		return c.Status(500).
			JSON(fiber.Map{"error": "Password requires at least 8 characters"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Failed to Hash", "details": err.Error()})
	}
	user.Password = string(hashedPassword)

	if err := db.DB.Create(user).Error; err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Failed to Create User", "details": err.Error()})
	}
	socketID, err := CreateWebSocketConnection(user.ID)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Failed to create WebSocket", "details": err.Error()})
	}
	token, exp, err := Create_JWT_Token(*user, socketID)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Failed to Create User", "details": err.Error()})
	}

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
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		Send_Error(ws, "User not found", err)
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

	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
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
	if err := db.DB.Where("email = ? AND id != ?", user.Email, userID).First(&existingUser).Error; err == nil {
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

	if err := db.DB.Save(user).Error; err != nil {

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

	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		Send_Error(ws, "User not found", err)
		return
	}

	if err := db.DB.Delete(user).Error; err != nil {
		Send_Error(ws, "Failed to update user", err)
		return
	}
	response := map[string]string{
		"Success": "Deleted User",
	}
	responseJson, _ := json.Marshal(response)
	Send_Message(ws, string(responseJson))
}
