package cli

import (
	"fmt"
)

// vmRun executes the command 'vm'.
func vmRun(args []string) int {
	if len(args) == 0 || args[0] == "help" {
		vmHelp()
		return 0
	}

	cmds := map[string]func([]string) int{
		"list":   func(a []string) int { return vmList(a) },
		"start":  func(a []string) int { return vmStart(a) },
		"status": func(a []string) int { return vmStatus(a) },
		"stop":   func(a []string) int { return vmStop(a) },
	}

	cmd, ok := cmds[args[0]]
	if !ok {
		fmt.Println("No such command: vm", args[0])
		fmt.Printf("\nFor usage information type:\n\n    rcstate vm help\n\n")
		return 1
	}

	return cmd(args)
}
