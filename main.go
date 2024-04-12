package main

import (
	"os"

	"github.com/OZahed/go-htmx/cmd"
)

func main() {
	if len(os.Args) < 2 {
		(cmd.HelpCmd{}).Execute(nil)
		os.Exit(1)
	}
	cmd.Execute(os.Args[1:])
}
