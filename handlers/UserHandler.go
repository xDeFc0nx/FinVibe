package handlers

import (
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func CreateUser(c *fiber.Ctx) error {
	user := new(types.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Unable to parse request body"})
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	user.Password = string(hashedPassword)

	token, exp, err := Create_JWT_Token(*user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create JWT token"})

	}

	if err := db.DB.Create(user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to Create User", "details": err.Error()})
	}

	socketID, err := CreateWebSocketConnection(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create WebSocket",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":  "User created successfully",
		"userID":   user.ID,
		"socketID": socketID,
		"token":    token,
		"exp":      exp,
	})
}

func Login_func(c *fiber.Ctx) error {

	user := new(types.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var foundUser types.User
	err := db.DB.Where("Email = ?", user.Email).First(&foundUser).Error
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid Email"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Wrong password"})
	}
	token, exp, err := Create_JWT_Token(foundUser)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create JWT token"})
	}

	cookie := fiber.Cookie{
		Name:     "jwt-token",
		Value:    token,
		Expires:  time.Unix(exp, 0),
		HTTPOnly: true,
	}
	socket := new(types.WebSocketConnection)
	err = db.DB.Where("user_id = ?", foundUser.ID).First(socket).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "WebSocket connection not found"})
	}

	socket.IsActive = true

	if err := db.DB.Save(socket).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update WebSocket connection"})
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{"message": "Success"})

}
func Logout_func(c *fiber.Ctx) error {
	// Extract JWT token from the cookie
	token := c.Cookies("jwt-token")
	if token == "" {
		return c.Status(401).JSON(fiber.Map{"error": "No JWT token provided"})
	}

	// Decode the JWT token and get the user ID
	userID, err := DecodeJWTToken(token) // You'll need to implement DecodeJWTToken
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Query the WebSocket connection using the userID
	socket := new(types.WebSocketConnection)
	err = db.DB.Where("user_id = ?", userID).First(socket).Error
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "WebSocket connection not found"})
	}

	// Set the WebSocket connection as inactive
	socket.IsActive = false
	if err := db.DB.Save(socket).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update WebSocket connection"})
	}

	// Clear the JWT cookie to log the user out
	cookie := fiber.Cookie{
		Name:     "jwt-token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{"message": "Success"})
}
