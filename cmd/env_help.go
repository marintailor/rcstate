package cmd

import "fmt"

// envHelp will print the usage information for env command
func envHelp() int {
	text := `
env command usage:
  rcstate env <argument> [option...]

Arguments:
  help    show usage

Options:
`
	fmt.Println(text)

	return 0
}
