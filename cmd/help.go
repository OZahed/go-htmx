/*
Using Viper or other packages for this spcefic situation is an overkill, we just want to seperate
Running server or Migrating database command and nothing else yet.
*/
package cmd

import (
	"fmt"

	"github.com/OZahed/go-htmx/internal/log"
)

func ExecuteHelp() {
	fmt.Printf(`
Subcommands:
	%shelp%s 	 shows help value
	%sserve%s 	 serving the server based on ENVs that are local to each pod/image
	%smigrate%s runs possible database migration files that are specified using DB_MIGRATION_PATH="..." on ENVs`+"\n\r\n\r",
		log.AnsiGreen, log.AnsiNoColor,
		log.AnsiGreen, log.AnsiNoColor,
		log.AnsiGreen, log.AnsiNoColor,
	)
}
