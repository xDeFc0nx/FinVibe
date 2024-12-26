package flag

import (
	"os"

	"github.com/xDeFc0nx/FinVibe/cmd/migrate"
)

func Flag() {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	arg := os.Args[1]

	if arg == "-m" {
		migrate.Migrate()
	}
}
