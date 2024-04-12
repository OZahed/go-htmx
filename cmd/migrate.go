package cmd

import (
	"fmt"
	"os"

	"github.com/OZahed/go-htmx/internal/log"
)

func ExecuteMigrate() {
	fmt.Fprintf(os.Stderr, "%s migrate command is%s %sNot Implemented %s",
		log.AnsiRed,
		log.AnsiNoColor,
		log.AnsiRedBG,
		log.AnsiNoColor,
	)

	ExecuteHelp()
}
