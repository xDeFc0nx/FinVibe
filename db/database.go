package db

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/xDeFc0nx/logger-go-pkg"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Conn() {
	if err := godotenv.Load(".env"); err != nil {
		logger.Error("%s", err.Error())
		return
	}

	var err error
	dsn := os.Getenv("DB_CONFIG")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {

		return
	}

	logger.Success("Successfully connected to Database!")
}
