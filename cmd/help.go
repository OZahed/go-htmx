/*
Using Viper or other packages for this spcefic situation is an overkill, we just want to seperate
Running server or Migrating database command and nothing else yet.
*/
package cmd

import (
	"fmt"
)

var _ Command = (*HelpCmd)(nil)

type HelpCmd struct{}

// Help implements Command.
func (HelpCmd) Help() HelpInfo {
	return HelpInfo{
		SubCmdName: "help",
		ShortDesc:  "displays help",
		Usage:      fmt.Sprintf("%s help", APP_NAME),
	}
}

// Name implements Command.
func (HelpCmd) Name() string {
	return "help"
}

func (HelpCmd) Execute(_ []string) {
	fmt.Println(APP_NAME)
	fmt.Println("cmd applications used for running an HTMX boiler plate application")
	fmt.Println("")
	fmt.Println("Sub Commands:")
	for name, cmd := range commands {
		help := cmd.Help()
		fmt.Printf("  %s:\n", name)
		fmt.Printf("\t%s\n", help.ShortDesc)
		if help.LongDesc != "" {
			fmt.Printf("\t%s\n", help.LongDesc)
		}
		fmt.Printf("\tUsage:\n\t%s\n", help.Usage)
	}
}
