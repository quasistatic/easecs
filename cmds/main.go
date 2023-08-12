package cmds

import (
	"strings"
)

const (
	HelpCommand = "help"
	ExecCommand = "exec"
	ListCommand = "list"
)

var aliasCmdToCmdNameMap = map[string]string{
	"help":   HelpCommand,
	"exec":   ExecCommand,
	"go":     ExecCommand,
	"jump":   ExecCommand,
	"enter":  ExecCommand,
	"launch": ExecCommand,
	"list":   ListCommand,
	"show":   ListCommand,
	"ls":     ListCommand,
	"all":    ListCommand,
}

type ICommand interface {
	Execute()
}

type BaseCmd struct {
	name string
	args []string
}

func IsCommandNameValid(cmdName string) bool {
	cmdName = strings.TrimSpace(cmdName)
	_, exists := aliasCmdToCmdNameMap[cmdName]
	return exists
}

func GetCommandInstance(cmdName string, args []string) ICommand {
	cmdName = strings.TrimSpace(cmdName)
	switch cmdName {
	case HelpCommand:
		return NewHelpCommand(args)
	case ListCommand:
		return NewListCommand(args)
	case ExecCommand:
		return NewExecCommand(args)
	}
	return nil
}
