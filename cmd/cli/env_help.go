package cli

import "fmt"

// envHelp will print the usage information for env command
func envHelp() {
	text := `
env command usage:
  rcstate env <command> [option...]

Commands:
  down    stop all resources in environment(s)
  help    show usage information
  show    show environment(a)
  up      start all resources in environment(s)

Options:
  -a, --all        show all environments
                   option is ignored when option "name" if provided

  --dry            Run the command without executing the logic

  -f, --format     Print the API request data of the command
                   Supported output formats: json

  -e, --env-file   environment file

  -h, --host       address of the remote host where the command will be executed

  -n, --name       environment name

Examples:
  Show all environments:

    rcstate env show \
      --all \
      -env-file <env_file>


  Show an environment:

    rcstate env show \
      --name <env_name> \
      --env-file <env_file>


  Show all environments or an environment with specific label(s):

    rcstate env show \
      { --all | --name <env_name> } \
      --label <label> \
      --env-file <env_file>


  Change state of all environments:

    rcstate env up/down \
      --all \
      --env-file <env_file>


  Change state of an environment from remote host:

    rcstate env up/down \
     --name <env_name> \
     --env-file <env_file> \
     --help <host_addr>


  Print the API request data in JSON format without executing the command:

    rcstate env show \
      --all \
      --env-file <env_file> \
      --format json \
      --dry
`
	fmt.Println(text)
}
