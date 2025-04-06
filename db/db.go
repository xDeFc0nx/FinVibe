package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/xDeFc0nx/logger-go-pkg"
)
var DB *pgx.Conn

func Conn() {
	DB, err := pgx.Connect(context.Background(), os.Getenv("DB_CONFIG"))
	if err != nil {
		logger.Error("Failed to connect to Database!", err)
		return
	}
	defer DB.Close(context.Background())

	logger.Success("Successfully connected to Database!")
}
