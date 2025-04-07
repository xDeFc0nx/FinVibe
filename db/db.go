package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func Conn() error {
	dbConfig := os.Getenv("DB_CONFIG")
	if dbConfig == "" {
		err := fmt.Errorf("DB_CONFIG environment variable not set")
		slog.Error("Database configuration error", slog.String("error", err.Error()))
		return err
	}

	var err error
	DB, err = pgx.Connect(context.Background(), dbConfig)
	if err != nil {
		slog.Error("Failed to connect to Database", slog.String("error", err.Error()))
		return fmt.Errorf("pgx.Connect failed: %w", err)
	}

	if err := DB.Ping(context.Background()); err != nil {
		slog.Error("Failed to ping database", slog.String("error", err.Error()))
		if DB != nil {
			_ = DB.Close(context.Background())
		}
		return fmt.Errorf("database ping failed: %w", err)
	}

	slog.Info("Successfully connected to Database!")
	return nil
}
