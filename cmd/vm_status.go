package cmd

import (
	"fmt"
)

// status return the status of the instance.
func (v VirtualMachine) status() int {
	if v.Opts.name == "" {
		fmt.Println("Please provide the instance's name.")
		return 1
	}

	status, err := v.Instances.Status(v.Opts.name)
	if err != nil {
		fmt.Printf("show status instance %q: %v\n", v.Opts.name, err)
		return 1
	}

	fmt.Println(status)

	return 0
}
