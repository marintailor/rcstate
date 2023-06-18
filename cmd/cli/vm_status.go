package cli

import (
	"fmt"

	"github.com/marintailor/rcstate/cmd/api/vm"
)

// vmStatus returns the status of the instance.
func vmStatus(args []string) int {
	v, err := vm.NewVirtualMachine(args)
	if err != nil {
		fmt.Println("list: new environment:", err)
		return 1
	}

	if v.Cfg.Name == "" {
		fmt.Println("Please provide the instance's name.")
		return 1
	}

	status, err := v.Instances.Status(v.Cfg.Name)
	if err != nil {
		fmt.Printf("show status instance %q: %v\n", v.Cfg.Name, err)
		return 1
	}

	fmt.Println(status)

	return 0
}
