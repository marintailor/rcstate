package cli

import (
	"encoding/json"
	"fmt"

	client "github.com/marintailor/rcstate/client/vm"
	"github.com/marintailor/rcstate/cmd/api/vm"
)

// list returns a table formatted list of instances.
func vmList(args []string) int {
	cfg := vm.Config{}

	if err := cfg.ParseFlags(args); err != nil {
		fmt.Println("get config:", err)
	}

	if cfg.Host != "" {
		return listRemote(&cfg)
	}

	return listLocal(&cfg)

}

// listLocal returns the list by executing the logic locally.
func listLocal(c *vm.Config) int {
	vm, err := vm.NewVirtualMachine(c.Project, c.Zone)
	if err != nil {
		fmt.Println("list: new environment:", err)
		return 1
	}

	list := vm.List()
	fmt.Println(list)

	return 0
}

// listRemote returns the list by sending a request to remote server.
func listRemote(c *vm.Config) int {
	j, err := json.Marshal(c)
	if err != nil {
		fmt.Println("marshal config:", err)
		return 1
	}

	if c.Format == "json" {
		fmt.Println(string(j))
	}

	if !c.Dry {
		data, err := client.List(string(j), c.Host)
		if err != nil {
			fmt.Println("client env down:", err)
			return 1
		}

		fmt.Println(data)
	}

	return 0
}
