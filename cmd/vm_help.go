package cmd

import "fmt"

// vmHelp prints the usage information for the command 'vm'.
func vmHelp() int {
	text := `
vm command usage:
  rcstate vm <argument> [option...]

Arguments:
  help      show usage
  list      list virtual machines

Options:
  -p, --project        Google Cloud Project ID

  -z, --zone           Google Cloud Zone name

Examples:

  List all instances in specific project and zone

    rcstate vm list \
      --project <project_name> \
      --zone <zone_name>
`
	fmt.Println(text)

	return 0
}
