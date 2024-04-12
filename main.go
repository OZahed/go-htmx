package main

import (
	"os"

	"github.com/OZahed/go-htmx/cmd"
)

func main() {
	if len(os.Args) < 2 {
		cmd.ExecuteHelp()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "serve":
		cmd.ExucuteServe()
	case "migrate":
		cmd.ExecuteMigrate()
	default:
		cmd.ExecuteHelp()
		if os.Args[1] != "help" {
			os.Exit(1)
		}
	}
}
