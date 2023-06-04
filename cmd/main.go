package cmd

import (
	"fmt"
)

// Run will execute the command with provided arguments.
func Run(args []string) int {
	if len(args) == 0 {
		help()
		return 0
	}

	return 0
}

// help shows the usage information.
func help() {
	text := `
Usage: rcstate <command> [args]

Commands:
`
	fmt.Println(text)
}
