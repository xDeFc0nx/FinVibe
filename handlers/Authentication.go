package handlers

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
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

func CreateJWTToken(
	userID string,
	connectionID string,
) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":       userID,
		"connection_id": connectionID,
		"exp":           exp,
	})

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
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				err := fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
				slog.Error("Unexpected signing method",
					slog.String("stack", fmt.Sprintf("%+v", errors.Wrap(err, ""))))

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
		slog.Error(MsgTokenDecodeFailed,
			slog.String("stack", fmt.Sprintf("%+v", errors.Wrap(err, ""))))
		return "", "", err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok &&
		parsedToken.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			slog.Error(
				MsgConnectionIDNotFound,
				slog.String("stack", fmt.Sprintf("%+v", errors.Wrap(err, ""))))
			return "", "", nil

		}

		connectionID, ok := claims["connection_id"].(string)
		if !ok {
			slog.Error(
				MsgConnectionIDNotFound,
				slog.String("stack", fmt.Sprintf("%+v", errors.Wrap(err, ""))))
			return "", "", nil

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
		JSendFail(c, data, fiber.StatusUnauthorized, err)
	}

	_, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		data := map[string]any{
			"message": "Invalid token format ",
		}
		JSendFail(c, data, fiber.StatusUnauthorized, err)
	}

	data := map[string]any{
		"message": "Authorized",
	}
	JSendSuccess(c, data)
	return nil
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
		JSendFail(c, data, fiber.StatusBadRequest, err)
		return err
	}
	var emailExists bool
	if err := db.DB.QueryRow(context.Background(), `
		SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)
		`, req.Email).Scan(&emailExists); err != nil {
		data := map[string]any{
			"message": "Failed to check email existence",
		}
		JSendError(c, data, fiber.StatusBadRequest, err)
	}
	if !emailExists {
		data := map[string]any{
			"message": "Email does not Exist",
		}
		JSendError(c, data, fiber.StatusBadRequest, nil)
		return nil
	}
	if err := db.DB.QueryRow(context.Background(), `
		SELECT connection_id 
		FROM web_sockets
		WHERE user_id = $1
		`, user.ID).Scan(&socket.ConnectionID); err != nil {
		data := map[string]any{
			"message": "failed to find ConnectionID",
		}
		JSendError(c, data, fiber.StatusNotFound, err)
		return err
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	); err != nil {
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
		JSendError(c, data, fiber.StatusNotFound, err)
		return err
	}
	token, exp, err := CreateJWTToken(user.ID, socket.ConnectionID)
	if err != nil {
		data := map[string]any{
			"message": "Failed to Create token",
		}
		JSendFail(c, data, fiber.StatusInternalServerError, err)
		return err

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
	JSendSuccess(c, data)
	return nil
}

func LogoutHandler(c *fiber.Ctx) error {
	token := c.Cookies(tokenName)
	c.Locals("csrf")
	if token == "" {
		data := map[string]any{
			"message": "Token is missing",
		}
		JSendFail(c, data, fiber.StatusBadRequest, nil)
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
	JSendSuccess(c, data)
	return nil
}
