package main

import (
	"os"

	"github.com/quasistatic/easecs/cmds"
)

func main() {
	// Clone the os.Args so that we can modify these as per case.
	mainArgs := os.Args

	// No arguments were passed, so assume "help" command.
	if len(mainArgs) == 1 {
		mainArgs = append(mainArgs, "help")
	}

	// The first useful argument is the core-command that needs to be triggered.
	// If that isn't a valid one, we over-write it as "help". The arguments are
	// updated to only include the faulty command name so that's printed for
	// user information.
	coreCommand := mainArgs[1]
	if !cmds.IsCommandNameValid(coreCommand) {
		coreCommand = cmds.HelpCommand
		var updatedArgs []string
		updatedArgs = append(updatedArgs, mainArgs[:1]...)
		updatedArgs = append(updatedArgs, coreCommand)
		updatedArgs = append(updatedArgs, mainArgs[1:]...)
		mainArgs = updatedArgs
	}

	// The relevant args are extracted flatly from the slice and passed as is.
	var cmdArgs []string
	if len(mainArgs) > 2 {
		cmdArgs = mainArgs[2:]
	}

	// Get the instance of the core command that is requested for.
	command := cmds.GetCommandInstance(coreCommand, cmdArgs)
	command.Execute()
}
