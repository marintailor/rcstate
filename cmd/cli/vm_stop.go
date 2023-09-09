package cli

import (
	"encoding/json"
	"fmt"

	client "github.com/marintailor/rcstate/client/vm"
	"github.com/marintailor/rcstate/cmd/api/vm"
)

// vmStop will stop an instance.
func vmStop(args []string) int {
	cfg := vm.Config{}

	if err := cfg.ParseFlags(args); err != nil {
		fmt.Println("get config:", err)
	}

	if cfg.Name == "" {
		fmt.Println("Please provide the instance's name.")
		return 1
	}

	if cfg.Host != "" {
		return stopRemote(&cfg)
	}

	return stopLocal(&cfg)

}

// stopLocal will stop an instance by executing the logic locally.
func stopLocal(c *vm.Config) int {
	v, err := vm.NewVirtualMachine(c.Project, c.Zone)
	if err != nil {
		fmt.Println("stop: new environment:", err)
		return 1
	}

	if err := v.Stop(c.Name); err != nil {
		fmt.Println("stop: ", err)
	}

	if c.DNS.RecordName != "" {
		dnsRecord := fmt.Sprintf("%s.%s", c.DNS.RecordName, c.DNS.Domain)
		c.Record(dnsRecord)
	}

	if c.Script.CMD != "" {
		c.ExecuteScript()
	}

	return 0
}

// stopRemote will stop an instance by sending a request to remote server.
func stopRemote(c *vm.Config) int {
	j, err := json.Marshal(c)
	if err != nil {
		fmt.Println("marshal config:", err)
		return 1
	}

	if c.Format == "json" {
		fmt.Println(string(j))
	}

	if !c.Dry {
		_, err = client.Stop(string(j), c.Host)
		if err != nil {
			fmt.Println("client env down:", err)
			return 1
		}
	}

	return 0
}
