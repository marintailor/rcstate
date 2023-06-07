package cmd

import "fmt"

// envHelp will print the usage information for env command
func envHelp() int {
	text := `
env command usage:
  rcstate env <argument> [option...]

Arguments:
  help    show usage
  show    show environment

Options:
  -a, --all        show all environments
                   option is ignored when option "name" if provided

  -e, -env-file    environment file

  -n, -name        environment name

Examples:
  Show all environments:

    rcstate env show -a -e <env_file>


  Show an environment:

    rcstate env show -n <env_name> -e <env_file>


  Show all environments or an environment with specific label(s):

    rcstate env show {-a | -n <env_name> } -l dev -e <env_file>
`
	fmt.Println(text)

	return 0
}
