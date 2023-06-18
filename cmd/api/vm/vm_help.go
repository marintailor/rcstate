package vm

import "fmt"

// vmHelp prints the usage information for the command 'vm'.
func vmHelp() int {
	text := `
vm command usage:
  rcstate vm <command> [option...]

Commands:
  help      show usage
  list      list virtual machines
  start     start the virtual machine
  status    show status of the virtual machine
  stop      stop the virtual machine

Options:
  -d, --domain         Domain for DNS record

  --dns-record-name    The DNS record name

  --dns-record-type    The DNS record type

  --external-ip        Use External IP address of instance for DNS record

  --ip                 Provide IP address for DNS record
                       Multiple addresses can be provided with comma delimiter

  -n, --name           Virtual Machine name

  -p, --project        Google Cloud Project ID

  -s, --script         Run shell command(s) on virtual machine with SSH connection
                       NOTE: command(s) must be wrapped in double quotes

  --ssh-key            Path to private key for SSH connection

  --ssh-port           Port number for SSH connection

  --ssh-user           Username for SSH connection

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


  Start an instance in specific project and zone, and create a DNS record

    rcstate vm start \
      --name <instance_name> \
      --project <project_name> \
      --zone <zone_name> \
      --domain <dns_domain> \
      --dns-record-name <record_name> \
      --dns-record-type <record_type>


  Show status of an instance in specific project and zone

    rcstate vm status \
      --name <instance_name> \
      --project <project_name> \
      --zone <zone_name>


  Stop an instance in specific project and zone

    rcstate vm stop \
      --name <instance_name> \
      --project <project_name> \
      --zone <zone_name>


  Start an instance and run shell commands AFTER the instance is started

    rcstate vm start \
      --name <instance_name> \
      --project <project_name> \
      --zone <zone_name> \
      --script "echo TEST > test-file" \
      --ip <ip_addr> \
      --ssh-key <path_to_key> \
      --ssh-port <port_number> \
      --ssh-user <username>


  Stop an instance and run shell commands BEFORE the instance is stopped

    rcstate vm stop \
      --name <instance_name> \
      --project <project_name> \
      --zone <zone_name> \
      --script "echo TEST > test-file" \
      --ip <ip_addr> \
      --ssh-key <path_to_key> \
      --ssh-port <port_number> \
      --ssh-user <username>
`
	fmt.Println(text)

	return 0
}
