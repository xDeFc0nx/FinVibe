package handlers

import (
	"encoding/json"
	"regexp"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xDeFc0nx/logger-go-pkg"
	"golang.org/x/crypto/bcrypt"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func CreateUser(c *fiber.Ctx) error {
	user := new(types.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Unable to parse request body", "details": err.Error()})
	}
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(user.Email) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email address", "email": user.Email})
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
	if user.Country == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Country is required"})
	}

	if len(user.Password) < 8 {
		return c.Status(500).JSON(fiber.Map{"error": "Password requires at least 8 characters"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to Hash", "details": err.Error()})
	}
	user.Password = string(hashedPassword)

	token, exp, err := Create_JWT_Token(*user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to Create User", "details": err.Error()})
	}

	if err := db.DB.Create(user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to Create User", "details": err.Error()})
	}

	socketID, err := CreateWebSocketConnection(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create WebSocket", "details": err.Error()})
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
		if err := ws.WriteMessage(websocket.TextMessage, []byte(`{"error":"User not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
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
	if err := ws.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}
}

func UpdateUser(ws *websocket.Conn, data json.RawMessage, userID string) {

	user := new(types.User)

	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if err := ws.WriteMessage(websocket.TextMessage, []byte(`{"error":"User not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	if err := json.Unmarshal(data, &user); err != nil {
		if err := ws.WriteMessage(websocket.TextMessage, []byte(`{"error":"Invalid data"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	if err := db.DB.Save(user).Error; err != nil {
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
	if err := ws.WriteMessage(websocket.TextMessage, responseData); err != nil {
		logger.Error("%s", err.Error())
	}
}
func DeleteUser(ws *websocket.Conn, data json.RawMessage, userID string) {

	user := new(types.User)

	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if err := ws.WriteMessage(websocket.TextMessage, []byte(`{"error":"User not found"}`+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	if err := db.DB.Delete(user).Error; err != nil {
		if err := ws.WriteMessage(websocket.TextMessage, []byte(`{"error":"Failed to Delete User"} `+err.Error())); err != nil {
			logger.Error("%s", err.Error())
		}
		return
	}

	if err := ws.WriteMessage(websocket.TextMessage, []byte(`{"Success": "User Deleted"}`)); err != nil {
		logger.Error("%s", err.Error())
	}

}
