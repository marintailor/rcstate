package cmd

import (
	"fmt"
)

// stop will stop an instance.
func (v VirtualMachine) stop() int {
	if v.Opts.name == "" {
		fmt.Println("Please provide the instance's name.")
		return 1
	}

	if v.Opts.script.cmd != "" {
		fmt.Println("=== Execute shell script")
		v.script()
	}

	fmt.Printf("=== Stop instance %q\n", v.Opts.name)
	if err := v.Instances.Stop(v.Opts.name); err != nil {
		fmt.Printf("stop instance %q: %v\n", v.Opts.name, err)
		return 1
	}

	fmt.Printf("=== Done\n\n")

	return 0
}
