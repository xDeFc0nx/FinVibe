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
