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
  start     start the virtual machine

Options:
  -n, --name           virtual machine name

  -p, --project        Google Cloud Project ID

  -z, --zone           Google Cloud Zone name

Examples:

  List all instances in specific project and zone

    rcstate vm list \
      --project <project_name> \
      --zone <zone_name>


  Start an instance in specific project and zone

    rcstate vm start \
      --name <instance_name> \
      --project <project_name> \
      --zone <zone_name>
`
	fmt.Println(text)

	return 0
}
