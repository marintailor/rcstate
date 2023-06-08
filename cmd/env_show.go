package cmd

import (
	"fmt"
	"strings"

	"github.com/marintailor/rcstate/api/gce"
)

// down will show all resources in environments.
func (e *Environments) show() int {
	switch {
	case e.name != "":
		env, err := e.getEnvironment(e.name)
		if err != nil {
			fmt.Println(err)
			return 1
		}

		if !env.checkLabel(e.label) {
			fmt.Printf("environment %q is not labeled with %q\n", e.name, e.label)
			return 1
		}

		showEnvironment(env)

		return 0
	case e.all:
		var count int

		for _, env := range e.Envs {
			if !env.checkLabel(e.label) {
				continue
			}

			showEnvironment(env)

			count++
		}

		if len(e.Envs) > 0 && count == 0 {
			fmt.Printf("no environment is labeled with %q\n", e.label)
		}
	}

	return 0
}

// showEnvironment will show all resources in specific environment.
func showEnvironment(env Environment) {
	fmt.Printf("\n%s\nENVIRONMENT: %s\nLABEL: %s\n%s\n", strings.Repeat("=", 40), env.Name, env.Label, strings.Repeat("=", 40))
	for i, g := range env.Group {
		if i > 0 {
			fmt.Println(strings.Repeat("-", 40))
		}

		fmt.Printf("\nGROUP: %s\nPROJECT: %s\nZONE: %s\n\n", g.Name, g.Project, g.Zone)
		fmt.Printf("VIRTUAL MACHINES\n")

		pw := padWidth(g.Resource.VM.Instance)

		list, err := gce.NewInstances(g.Project, g.Zone).GetList()
		if err != nil {
			fmt.Printf("new instances: %v", err)
		}

		for j, instance := range g.Resource.VM.Instance {
			var status string
			for _, inst := range list {
				if inst.Name == instance.Name {
					status = inst.Status
				}
			}
			padding := strings.Repeat(" ", pw-len(instance.Name))
			fmt.Printf("%d. %s%sStatus: %s\n", j+1, instance.Name, padding, status)
		}
		fmt.Println()
	}
}

// padWidth returns width of the pad.
func padWidth(instances []Instance) int {
	pw := 16
	for _, instance := range instances {
		if len(instance.Name) >= pw {
			pw = len(instance.Name) + 2
		}
	}

	return pw
}
