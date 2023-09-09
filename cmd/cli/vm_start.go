package cli

import (
	"encoding/json"
	"fmt"

	client "github.com/marintailor/rcstate/client/vm"
	"github.com/marintailor/rcstate/cmd/api/vm"
)

// vmStart will start an instance.
func vmStart(args []string) int {
	cfg := vm.Config{}

	if err := cfg.ParseFlags(args); err != nil {
		fmt.Println("get config:", err)
	}

	if cfg.Name == "" {
		fmt.Println("Please provide the instance's name.")
		return 1
	}

	if cfg.Host != "" {
		return startRemote(&cfg)
	}

	return startLocal(&cfg)

}

// startLocal will start an instance by executing the logic locally.
func startLocal(c *vm.Config) int {
	v, err := vm.NewVirtualMachine(c.Project, c.Zone)
	if err != nil {
		fmt.Println("start: new environment:", err)
		return 1
	}

	if err := v.Start(c.Name); err != nil {
		fmt.Println("start: ", err)
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

// startRemote will start an instance by sending a request to remote server.
func startRemote(c *vm.Config) int {
	j, err := json.Marshal(c)
	if err != nil {
		fmt.Println("marshal config:", err)
		return 1
	}

	if c.Format == "json" {
		fmt.Println(string(j))
	}

	if !c.Dry {
		_, err = client.Start(string(j), c.Host)
		if err != nil {
			fmt.Println("client env down:", err)
			return 1
		}
	}

	return 0
}
