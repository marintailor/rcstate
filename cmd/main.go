package cmd

import (
	"fmt"

	"github.com/marintailor/rcstate/cmd/cli"
)

// Run executes the specified command.
func Run(args []string) int {
	if len(args) == 0 {
		help()
		return 0
	}

	return cli.Run(args)
}

// TODO: update
// help shows the usage information.
func help() {
	text := `
Usage: rcstate <command> [options...]

Commands:
  env     manage declared environments
  vm      manage state of virtual machine instance
`
	fmt.Println(text)
}
