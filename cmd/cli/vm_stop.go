package cli

import (
	"fmt"

	"github.com/marintailor/rcstate/cmd/api/vm"
)

// vmStop will stop an instance.
func vmStop(args []string) int {
	v, err := vm.NewVirtualMachine(args)
	if err != nil {
		fmt.Println("list: new environment:", err)
		return 1
	}

	if v.Cfg.Name == "" {
		fmt.Println("Please provide the instance's name.")
		return 1
	}

	if v.Cfg.Script.CMD != "" {
		v.Cfg.ExecuteScript()
	}

	if err := v.Instances.Stop(v.Cfg.Name); err != nil {
		fmt.Printf("stop instance %q: %v\n", v.Cfg.Name, err)
		return 1
	}

	return 0
}
