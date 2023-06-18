package cli

import (
	"fmt"

	"github.com/marintailor/rcstate/cmd/api/vm"
)

// vmStart will start an instance.
func vmStart(args []string) int {
	v, err := vm.NewVirtualMachine(args)
	if err != nil {
		fmt.Println("list: new environment:", err)
		return 1
	}

	if v.Cfg.Name == "" {
		fmt.Println("Please provide the instance's name.")
		return 1
	}

	if err := v.Instances.Start(v.Cfg.Name); err != nil {
		fmt.Printf("Could not start instance %q: %v\n", v.Cfg.Name, err)
		return 1
	}

	if v.Cfg.DNS.RecordName != "" {
		dnsRecord := fmt.Sprintf("%s.%s", v.Cfg.DNS.RecordName, v.Cfg.DNS.Domain)
		v.Cfg.Record(dnsRecord)
	}

	if v.Cfg.Script.CMD != "" {
		v.Cfg.ExecuteScript()
	}

	return 0
}
