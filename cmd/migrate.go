package cmd

import (
	"fmt"
	"os"

	"github.com/OZahed/go-htmx/internal/logger"
)

var _ Command = (*MigratCmd)(nil)

type MigratCmd struct{}

// Help implements Command.
func (MigratCmd) Help() HelpInfo {
	return HelpInfo{
		SubCmdName: "migrate",
		Usage:      fmt.Sprintf("%s migrate --path=$MIGRATE_PATH [opt...]", APP_NAME),
		ShortDesc:  "database migration, migration files should be ",
	}
}

// Name implements Command.
func (MigratCmd) Name() string {
	return "migrate"
}

func (MigratCmd) Execute(_ []string) {
	fmt.Fprintf(os.Stderr, "%s migrate command is%s %sNot Implemented %s",
		logger.AnsiRed,
		logger.AnsiNoColor,
		logger.AnsiRedBG,
		logger.AnsiNoColor,
	)
}
