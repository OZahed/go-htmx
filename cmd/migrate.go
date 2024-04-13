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
		LongDesc:   `Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.`,
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
