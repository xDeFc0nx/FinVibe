package handlers

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/xDeFc0nx/logger-go-pkg"

	"github.com/xDeFc0nx/FinVibe/types"
)

func Create_JWT_Token(user types.User) (string, int64, error) {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("Error loading.env file")
		os.Exit(1)
	}

	SECRET_KEY := os.Getenv("SECRET_KEY")
	exp := time.Now().Add(time.Minute * 30).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = exp
	t, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", 0, err
	}
	return t, exp, nil
}
func DecodeJWTToken(token string) (string, error) {
	// Remove the "Bearer " prefix if it's there
	token = strings.TrimPrefix(token, "Bearer ")

	// Parse the token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method matches the expected method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("your-secret-key"), nil // Replace with your actual secret key
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

	return c.Status(200).JSON("Authorized")
}
