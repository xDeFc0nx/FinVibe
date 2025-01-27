package handlers

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

var tokenName = "jwt-token"

func Create_JWT_Token(
	user types.User,
	connectionID string,
) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 30).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["connectionID"] = connectionID
	claims["exp"] = exp
	t, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", 0, err
	}
	return t, exp, nil
}

func DecodeJWTToken(token string) (string, string, error) {
	token = strings.TrimPrefix(
		token,
		"Bearer ",
	)

	parsedToken, err := jwt.Parse(
		token,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(
				os.Getenv("SECRET_KEY"),
			), nil
		},
	)
	if err != nil {
		slog.Error("Failed to parse token", slog.String("error", err.Error()))
		return "", "", err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok &&
		parsedToken.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			slog.Error(
				"userID not found in claims",
			)
			return "", "", err

		}

		connectionID, ok := claims["connectionID"].(string)
		if !ok {
			slog.Error(
				"connectionID not found in claims",
			)
			return "", "", err

		}

		return userID, connectionID, nil
	}
	slog.Error(
		"invalid token or claims",
	)
	return "", "", err
}

func CheckAuth(ws *fiber.Ctx) error {
	cookie := ws.Cookies(tokenName)

	token, err := jwt.ParseWithClaims(
		cookie,
		&jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)
	if err != nil || !token.Valid {
		return ws.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	_, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return ws.Status(401).JSON(fiber.Map{"error": "Invalid token format"})
	}

	return ws.Status(200).JSON(fiber.Map{"message": "Authorized"})
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

	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.Password),
		[]byte(user.Password),
	)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Wrong password"})
	}

	socket := new(types.WebSocketConnection)
	err = db.DB.Where("user_id = ?", foundUser.ID).First(socket).Error
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "WebSocket connection not found"})
	}
	token, exp, err := Create_JWT_Token(foundUser, socket.ConnectionID)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Failed to create JWT token"})
	}
	cookie := fiber.Cookie{
		Name:     tokenName,
		Value:    token,
		Expires:  time.Unix(exp, 0),
		HTTPOnly: true,
	}

	socket.IsActive = true

	if err := db.DB.Save(socket).Error; err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Failed to update WebSocket connection"})
	}
	ConnectionID := socket.ConnectionID
	c.Cookie(&cookie)
	return c.JSON(
		fiber.Map{
			"message":       "Success",
			"Connection ID": ConnectionID,
			"token":         cookie.Value,
		},
	)
}

func LogoutHandler(c *fiber.Ctx) error {
	socket := new(types.WebSocketConnection)

	token := c.Cookies(tokenName)

	if token == "" {
		log.Println("Token is missing")
		return c.Status(400).JSON(fiber.Map{"error": "Token is missing"})
	}
	userID, _, err := DecodeJWTToken(token)
	if err != nil {
		log.Printf("Failed to decode JWT token: %v\n", err)
	}

	if err := db.DB.Where("user_id = ?", userID).First(socket).Error; err != nil {
		log.Println("Failed to find WebSocket connection: ", err)
		return c.Status(500).
			JSON(fiber.Map{"error": "Failed to find WebSocket connection"})
	}

	socket.IsActive = false
	if err := db.DB.Save(socket).Error; err != nil {
		log.Println("Failed to update WebSocket connection: ", err)
		return c.Status(500).
			JSON(fiber.Map{"error": "Failed to update WebSocket connection"})
	}

	cookie := fiber.Cookie{
		Name:     tokenName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}
