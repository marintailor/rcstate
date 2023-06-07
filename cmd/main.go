package cmd

import (
	"fmt"
)

// Run executes the specified command.
func Run(args []string) int {
	if len(args) == 0 {
		help()
		return 0
	}

	cmds := map[string]func([]string) int{
		"env": func(a []string) int { return envRun(a) },
		"vm":  func(a []string) int { return vmRun(a) },
	}

	cmd, ok := cmds[args[0]]
	if !ok {
		fmt.Println("No such cmd:", args[0])
		help()
		return 1
	}

	return cmd(args[1:])
}

// help shows the usage information.
func help() {
	text := `
Usage: rcstate <command> [args]

Commands:
  vm      manage state of virtual machine instances
`
	fmt.Println(text)
}
