package cmd

import (
	"fmt"
	"runtime"

	"github.com/OZahed/go-htmx/internal/logger"
)

var (
	APP_VERSION string
	APP_NAME    string
	GIT_HEAD    string
	BUILD_AT    string
)

var _ Command = (*VersionCmd)(nil)

type VersionCmd struct{}

// Help implements Command.
func (VersionCmd) Help() HelpInfo {
	return HelpInfo{
		SubCmdName: "version",
		Usage:      fmt.Sprintf("%s version", APP_NAME),
		ShortDesc:  "displays version and build information",
	}
}

// Name implements Command.
func (v VersionCmd) Name() string {
	return "version"
}

func (VersionCmd) Execute(_ []string) {
	fmt.Printf(`
Version Informaton
----------------------
Application: %s 
app version: %s 
head hash:   %s 
build time:  %s 
Go version:  %s 
----------------------
`,

		logger.AnsiYellow+APP_NAME+logger.AnsiNoColor,
		logger.AnsiYellow+APP_VERSION+logger.AnsiNoColor,
		logger.AnsiYellow+GIT_HEAD+logger.AnsiNoColor,
		logger.AnsiYellow+BUILD_AT+logger.AnsiNoColor,
		logger.AnsiYellow+runtime.Version()+logger.AnsiNoColor,
	)
}
