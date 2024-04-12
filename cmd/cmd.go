package cmd

var commands Commands

func init() {
	commands = make(Commands)

	commands.Register(ServeCmd{})
	commands.Register(HelpCmd{})
	commands.Register(MigratCmd{})
	commands.Register(VersionCmd{})
}

type Commands map[string]Command

func (c Commands) Register(cmd Command) {
	c[cmd.Name()] = cmd
}

type Command interface {
	Execute(args []string)
	Help() HelpInfo
	Name() string
}

type HelpInfo struct {
	SubCmdName, Usage, ShortDesc, LongDesc string
}

func Execute(args []string) {
	if len(args) == 0 {
		return
	}

	cmd, ok := commands[args[0]]
	if !ok {
		cmd = HelpCmd{}
	}

	cmd.Execute(args)
}
