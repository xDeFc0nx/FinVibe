package migrate

import (
	"github.com/xDeFc0nx/logger-go-pkg"

	"github.com/xDeFc0nx/FinVibe/db"
	"github.com/xDeFc0nx/FinVibe/types"
)

func Migrate() {
	if err := db.DB.AutoMigrate(

		&types.WebSocketConnection{},
		&types.User{},
		&types.Transaction{},
		&types.Recurring{},
	); err != nil {
		return
	}
	logger.Success("Successfully Migrated!")

}
