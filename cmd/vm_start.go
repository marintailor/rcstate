package cmd

import (
	"fmt"
)

// start will start an instance.
func (v VirtualMachine) start() int {
	if v.Opts.name == "" {
		fmt.Println("Please provide the instance's name.")
		return 1
	}

	fmt.Printf("=== Start instance %q\n", v.Opts.name)
	if err := v.Instances.Start(v.Opts.name); err != nil {
		fmt.Printf("Could not start instance %q: %v\n", v.Opts.name, err)
		return 1
	}

	fmt.Printf("=== Done\n\n")

	return 0
}
