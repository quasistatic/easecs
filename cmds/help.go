package cmds

import (
	"fmt"
	"strings"
)

type Help BaseCmd

func NewHelpCommand(args []string) *Help {
	return &Help{
		name: HelpCommand,
		args: args,
	}
}

func (d *Help) Execute() {
	helpText := `
EaseCS Version 0.0.1

Usage:
	easecs <command> <args...>

Commands:
	easecs help
	easecs list
	easecs exec <cluster> <service> <container>`

	if len(d.args) > 0 {
		fmt.Println("[ERROR] Unknown command '", strings.Join(d.args, " "))
	}
	fmt.Println(helpText)
}
