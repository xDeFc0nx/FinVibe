package handlers

import (
	"context"
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
		func(token *jwt.Token) (any, error) {
			if _, err := token.Method.(*jwt.SigningMethodHMAC); !err {
				slog.Error(
					"Unexpected signing method",
					"method", token.Method.Alg(),
				)
				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					token.Method.Alg(),
				)
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

func CheckAuth(c *fiber.Ctx) error {
	cookie := c.Cookies(tokenName)

	token, err := jwt.ParseWithClaims(
		cookie,
		&jwt.MapClaims{},
		func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)
	if err != nil || !token.Valid {
		data := map[string]any{
			"message": "Token is not valid",
		}
		return JSendFail(c, data, fiber.StatusUnauthorized)
	}

	_, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		data := map[string]any{
			"message": "Invalid token format ",
		}
		return JSendFail(c, data, fiber.StatusUnauthorized)
	}

	data := map[string]any{
		"message": "Authorized",
	}
	return JSendSuccess(c, data)
}

func LoginHandler(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Request
	user := new(types.User)
	socket := new(types.WebSocket)
	c.Locals("csrf")
	if err := c.BodyParser(&req); err != nil {
		data := map[string]any{
			"Message": "invalid body",
		}
		return JSendFail(c, data, fiber.StatusBadRequest)
	}
	var foundUser types.User
	err := db.DB.QueryRow(context.Background(), `
		SELECT  email, 
		FROM users
		WHERE email = $1
	`, req.Email).Scan(
		&user.Email,
	)
	if err != nil {
		slog.Error(
			"Email Does not exist",
			"error",
			err,
			"IP",
			c.IP(),
		)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		slog.Error(
			"Password Login Attempt",
			"Email",
			user.Email,
			"error",
			err,
			"IP",
			c.IP(),
		)

		data := map[string]any{
			"error": "Wrong Email or password",
		}
		return JSendError(c, data, fiber.StatusNotFound)
	}
	token, exp, err := Create_JWT_Token(foundUser, socket.ConnectionID)
	if err != nil {
		data := map[string]any{
			"message": "Failed to Create token",
		}
		return JSendFail(c, data, fiber.StatusInternalServerError)

	}
	c.Cookie(&fiber.Cookie{
		Name:     tokenName,
		Value:    token,
		Expires:  time.Unix(exp, 0),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})
	slog.Info(
		"Authentication successful",
		"ID",
		user.ID,
		"IP",
		c.IP(),
	)
	data := map[string]any{
		"message": "Authorized",
		"expires": exp,
	}
	return JSendSuccess(c, data)
}

func LogoutHandler(c *fiber.Ctx) error {
	token := c.Cookies(tokenName)
	userID := c.Locals("userID").(string)
	c.Locals("csrf")
	if token == "" {
		data := map[string]any{
			"message": "Token is missing",
		}
		return JSendFail(c, data, fiber.StatusBadRequest)
	}
	userID, _, err := DecodeJWTToken(token)
	if err != nil {
		log.Printf(" %v\n", err)
		slog.Info(
			"Failed to decode JWT token",
			"ID",
			userID,
			"error",
			err,
			"IP",
			c.IP(),
		)

	}
	if _, err := db.DB.Exec(context.Background(), `
 UPDATE websockets 
 SET is_active = false
 WHERE user_id = $1
`, userID); err != nil {
		slog.Error(
			"Failed to update websocket",
			"ID",
			userID,
			"error",
			err,
			"IP",
			c.IP(),
		)
	}
	cookie := fiber.Cookie{
		Name:     tokenName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	}
	c.Cookie(&cookie)
	slog.Info(
		"Logged out successfully",
		"ID",
		userID,
		"IP",
		c.IP(),
	)
	data := map[string]any{
		"message": "Logged out",
	}
	return JSendSuccess(c, data)
}
