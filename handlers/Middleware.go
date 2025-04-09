package handlers

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(c *fiber.Ctx) error {
	tokenString := c.Cookies(tokenName)
	if tokenString == "" {
		slog.Error("Missing token", "IP", c.IP())
		data := map[string]interface{}{
			"message": "Unauthorized: Missing token",
		}

		JSendError(c, data, fiber.StatusUnauthorized, nil)
		return nil

	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)

	if err != nil || !token.Valid {
		slog.Error("Invalid token", "error", err, "IP", c.IP())
		data := map[string]interface{}{
			"message": "Unauthorized: Invalid token",
		}
		JSendError(c, data, fiber.StatusUnauthorized, nil)
		return nil
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		slog.Error("Invalid token claims", "IP", c.IP())
		data := map[string]interface{}{
			"message": "Unauthorized: Invalid token claims",
		}
		JSendError(c, data, fiber.StatusUnauthorized, nil)
		return nil
	}

	if claims.Subject == "" {
		slog.Error("Missing subject in token", "IP", c.IP())
		data := map[string]interface{}{
			"message": "Unauthorized: Missing subject in token",
		}
		JSendError(c, data, fiber.StatusUnauthorized, nil)
		return nil
	}

	c.Locals("userID", claims.Subject)
	return c.Next()
}
