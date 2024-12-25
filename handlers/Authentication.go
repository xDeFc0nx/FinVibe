package handlers

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/xDeFc0nx/logger-go-pkg"
	"golang.org/x/crypto/bcrypt"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func Create_JWT_Token(user types.User) (string, int64, error) {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("Error loading.env file")
		os.Exit(1)
	}

	exp := time.Now().Add(time.Minute * 30).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = exp
	t, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", 0, err
	}
	return t, exp, nil
}
func DecodeJWTToken(token string) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("Error loading.env file")
		os.Exit(1)
	}
	token = strings.TrimPrefix(token, "Bearer ")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method matches the expected method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil // Replace with your actual secret key
	})
	if err != nil {
		return "", fmt.Errorf("Failed to parse token: %v", err)
	}

	// Check if the token is valid and extract claims
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", fmt.Errorf("User ID not found in token")
		}
		return userID, nil
	}

	return "", fmt.Errorf("Invalid token or claims")
}
func CheckAuth(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt-token")

	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token format"})
	}

	userID := (*claims)["userID"].(string)

	return c.Status(200).JSON(fiber.Map{"message": "Authorized", "userID": userID})
}

func LoginHandler(c *fiber.Ctx) error {

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
	ConnectionID := socket.ConnectionID
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{"message": "Success", "Conection ID": ConnectionID})

}
func LogoutHandler(conn *websocket.Conn, userID string) {

	socket := new(types.WebSocketConnection)
	err := db.DB.Where("user_id = ?", userID).First(socket).Error
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"Failed to get websocket connection"}`+err.Error()))
	}

	socket.IsActive = false
	if err := db.DB.Save(socket).Error; err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"Failed to Update Socket Connection"}`+err.Error()))
	}

	conn.WriteMessage(websocket.TextMessage, []byte(`{"Success": "Logedout"}`))

	conn.Close()
}
