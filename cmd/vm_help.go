package cmd

import "fmt"

// vmHelp prints the usage information for the command 'vm'.
func vmHelp() int {
	text := `
vm command usage:
  rcstate vm <argument> [option...]

Arguments:
  help      show usage

Options:
`
	fmt.Println(text)

	return 0
}
