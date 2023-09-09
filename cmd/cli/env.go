package cli

import (
	"fmt"
)

// envRun runs the logic for the env command
func envRun(args []string) int {
	if len(args) == 0 || args[0] == "help" {
		envHelp()
		return 0
	}

	commands := map[string]func([]string) int{
		"down": func(a []string) int { return envDown(a) },
		"show": func(a []string) int { return envShow(a) },
		"up":   func(a []string) int { return envUp(a) },
	}

	command, ok := commands[args[0]]
	if !ok {
		fmt.Println("No such command: env", args[0])
		fmt.Printf("\nFor usage information type:\n\n    rcstate env help\n\n")
		return 1
	}

	return command(args)
}
