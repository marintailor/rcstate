package cli

import (
	"encoding/json"
	"fmt"

	client "github.com/marintailor/rcstate/client/vm"
	"github.com/marintailor/rcstate/cmd/api/vm"
)

// vmStatus return the status of the instance.
func vmStatus(args []string) int {
	cfg := vm.Config{}

	if err := cfg.ParseFlags(args); err != nil {
		fmt.Println("get config:", err)
	}

	if cfg.Name == "" {
		fmt.Println("Please provide the instance's name.")
		return 1
	}

	if cfg.Host != "" {
		return statusRemote(&cfg)
	}

	return statusLocal(&cfg)

}

// statusLocal return the status of the instance by executing the logic locally.
func statusLocal(c *vm.Config) int {
	v, err := vm.NewVirtualMachine(c.Project, c.Zone)
	if err != nil {
		fmt.Println("status: new environment:", err)
		return 1
	}

	status, err := v.Status(c.Name)
	if err != nil {
		fmt.Println("status: ", err)
	}

	fmt.Println(status)

	return 0
}

// statusRemote return the status of the instance by sending a request to remote server.
func statusRemote(c *vm.Config) int {
	j, err := json.Marshal(c)
	if err != nil {
		fmt.Println("marshal config:", err)
		return 1
	}

	if c.Format == "json" {
		fmt.Println(string(j))
	}

	if !c.Dry {
		data, err := client.Status(string(j), c.Host)
		if err != nil {
			fmt.Println("client env down:", err)
			return 1
		}

		fmt.Println(data)
	}

	return 0
}
