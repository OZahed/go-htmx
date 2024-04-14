package cmd

import (
	"fmt"
	"os"

	"github.com/OZahed/go-htmx/internal/logger"
)

var _ Command = (*MigrateCmd)(nil)

type MigrateCmd struct{}

// Help implements Command.
func (MigrateCmd) Help() HelpInfo {
	return HelpInfo{
		SubCmdName: "migrate",
		Usage:      fmt.Sprintf("%s migrate --path=$MIGRATE_PATH [opt...]", APP_NAME),
		ShortDesc:  "database migration, migration files should be ",
		LongDesc:   "",
	}
}

// Name implements Command.
func (MigrateCmd) Name() string {
	return "migrate"
}

func (MigrateCmd) Execute(_ []string) {
	// Read Files in a directory, check for their crate time, check last migration time on the DB, execute all the
	// SQL files that are generated after last migration
	fmt.Fprintf(os.Stderr, "%s migrate command is%s %sNot Implemented %s",
		logger.AnsiRed,
		logger.AnsiNoColor,
		logger.AnsiRedBG,
		logger.AnsiNoColor,
	)

	fmt.Println("")
	HelpCmd{}.Execute(nil)
}
