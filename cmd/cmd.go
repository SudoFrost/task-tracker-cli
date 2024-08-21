package cmd

type Command struct {
	name    string
	handler func([]string)
}

func NewCommand(name string, handler func([]string)) Command {
	return Command{name, handler}
}

type CLI struct {
	args     []string
	commands []Command
}

func NewCLI(args []string) *CLI {
	return &CLI{args: args}
}

func (cli *CLI) AddCommand(cmd Command) {
	cli.commands = append(cli.commands, cmd)
}

func (cli *CLI) Run() {
	if len(cli.args) == 0 {
		return
	}
	for _, cmd := range cli.commands {
		if cmd.name == cli.args[0] {
			args := []string{}
			if len(cli.args) > 1 {
				args = cli.args[1:]
			}
			cmd.handler(args)
		}
	}
}
