package cli

import (
	"flag"
	"fmt"
)

type config struct {
	serverPort string
}

func Run(args []string) int {
	var cfg = config{}
	cfg.getConfig(args)

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

func (c *config) getConfig(args []string) {
	f := flag.NewFlagSet(args[0], flag.ContinueOnError)

	f.StringVar(&c.serverPort, "server", "", "Run in server mode")
	f.StringVar(&c.serverPort, "s", "", "Run in server mode")

	if err := f.Parse(args); err != nil {
		fmt.Println(err)
	}
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
