package flag

import (
	"os"

	"github.com/xDeFc0nx/FinVibe/cmd/migrate"
)

func Flag() {
	defer func() {
		if r := recover(); r != nil {
			// If there's an out of range error, just ignore it and continue
		}
	}()
	arg := os.Args[1]

	if arg == "-m" {
		migrate.Migrate()
	} else {

	}
}
